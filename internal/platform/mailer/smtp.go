package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	From       string
	Encryption string
	Timeout    time.Duration
	RootCAs    *x509.CertPool
}

type SMTP struct {
	config Config
	from   *mail.Address
}

func NewSMTP(config Config) (*SMTP, error) {
	from, err := mail.ParseAddress(config.From)
	if err != nil || strings.ContainsAny(config.From, "\r\n") || config.Host == "" || strings.ContainsAny(config.Host, "\r\n") ||
		config.Port < 1 || config.Port > 65535 || config.Username == "" || config.Password == "" || config.Timeout <= 0 ||
		(config.Encryption != "starttls" && config.Encryption != "tls") {
		return nil, errors.New("SMTP requires a valid host, port, username, password, sender, TLS mode and timeout")
	}
	return &SMTP{config: config, from: from}, nil
}

func (s *SMTP) Send(ctx context.Context, to, subject, body string) error {
	recipient, err := mail.ParseAddress(to)
	if err != nil || recipient.Address != to || strings.ContainsAny(to+subject, "\r\n") {
		return errors.New("invalid email recipient or subject")
	}
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	address := net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port))
	connection, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("SMTP connect: %w", err)
	}
	defer connection.Close()
	stop := context.AfterFunc(ctx, func() { _ = connection.Close() })
	defer stop()
	if deadline, ok := ctx.Deadline(); ok {
		if err := connection.SetDeadline(deadline); err != nil {
			return err
		}
	}
	tlsConfig := &tls.Config{ServerName: s.config.Host, MinVersion: tls.VersionTLS12, RootCAs: s.config.RootCAs}
	var transport net.Conn = connection
	if s.config.Encryption == "tls" {
		encrypted := tls.Client(connection, tlsConfig)
		if err := encrypted.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("SMTP TLS: %w", err)
		}
		transport = encrypted
	}
	client, err := smtp.NewClient(transport, s.config.Host)
	if err != nil {
		return fmt.Errorf("SMTP greeting: %w", err)
	}
	defer client.Close()
	if s.config.Encryption == "starttls" {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return errors.New("SMTP server does not support required STARTTLS")
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("SMTP STARTTLS: %w", err)
		}
	}
	if err := client.Auth(smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)); err != nil {
		return fmt.Errorf("SMTP authentication: %w", err)
	}
	if err := client.Mail(s.from.Address); err != nil {
		return fmt.Errorf("SMTP sender: %w", err)
	}
	if err := client.Rcpt(recipient.Address); err != nil {
		return fmt.Errorf("SMTP recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("SMTP data: %w", err)
	}
	message := s.message(recipient, subject, body)
	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("SMTP write: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("SMTP acceptance: %w", err)
	}
	_ = client.Quit()
	return nil
}

func (s *SMTP) message(to *mail.Address, subject, body string) []byte {
	var message bytes.Buffer
	fmt.Fprintf(&message, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n", s.from.String(), to.String(), mime.QEncoding.Encode("utf-8", subject), time.Now().UTC().Format(time.RFC1123Z))
	writer := quotedprintable.NewWriter(&message)
	_, _ = writer.Write([]byte(body))
	_ = writer.Close()
	return message.Bytes()
}
