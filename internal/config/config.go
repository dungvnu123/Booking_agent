package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net"
	"net/mail"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

type Config struct {
	Addr                 string
	Mode                 string
	DatabaseURL          string
	DBMaxConns           int32
	Location             *time.Location
	CORSOrigins          []string
	SolverURL            string
	SolverTimeout        time.Duration
	StorageDir           string
	JWTAccessSecret      string
	JWTRefreshSecret     string
	JWTIssuer            string
	JWTAccessTTL         time.Duration
	JWTRefreshTTL        time.Duration
	SMTPHost             string
	SMTPPort             int
	SMTPUsername         string
	SMTPPassword         string
	SMTPFrom             string
	SMTPEncryption       string
	SMTPTimeout          time.Duration
	EmailVerificationURL string
	EmailVerificationTTL time.Duration
	EmailResendCooldown  time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, errors.New("cannot parse .env")
	}
	return FromEnv()
}

// FromEnv is separate from dotenv loading for deterministic tests.
func FromEnv() (Config, error) {
	cfg := Config{
		Addr:             value("SERVER_ADDR", "127.0.0.1:8081"),
		Mode:             value("GIN_MODE", "debug"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		SolverURL:        value("ROSTERING_URL", "http://localhost:8000"),
		StorageDir:       value("STORAGE_DIR", ".data/uploads"),
		JWTAccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
		JWTRefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
		JWTIssuer:        value("JWT_ISSUER", "abclinic-backend"),
		SMTPHost:         os.Getenv("SMTP_HOST"), SMTPUsername: os.Getenv("SMTP_USERNAME"), SMTPPassword: os.Getenv("SMTP_PASSWORD"), SMTPFrom: os.Getenv("SMTP_FROM"),
		SMTPEncryption:       value("SMTP_ENCRYPTION", "starttls"),
		EmailVerificationURL: value("EMAIL_VERIFICATION_URL", "http://localhost:8081/api/auth/verify-email"),
	}
	if cfg.Mode != "debug" && cfg.Mode != "release" && cfg.Mode != "test" {
		return Config{}, errors.New("invalid GIN_MODE")
	}
	_, port, err := net.SplitHostPort(cfg.Addr)
	if err != nil {
		return Config{}, errors.New("SERVER_ADDR must contain host:port")
	}
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		return Config{}, errors.New("SERVER_ADDR port must be 1..65535")
	}
	cfg.Location, err = time.LoadLocation(value("APP_TIMEZONE", "Asia/Ho_Chi_Minh"))
	if err != nil {
		return Config{}, errors.New("invalid APP_TIMEZONE")
	}
	maxConns, err := strconv.Atoi(value("DB_MAX_CONNS", "10"))
	if err != nil || maxConns < 1 || maxConns > 100 {
		return Config{}, errors.New("DB_MAX_CONNS must be 1..100")
	}
	cfg.DBMaxConns = int32(maxConns)
	cfg.SolverTimeout, err = time.ParseDuration(value("ROSTERING_TIMEOUT", "30s"))
	if err != nil || cfg.SolverTimeout <= 0 {
		return Config{}, errors.New("invalid ROSTERING_TIMEOUT")
	}
	solverURL, err := url.Parse(cfg.SolverURL)
	if err != nil || solverURL.Host == "" || (solverURL.Scheme != "http" && solverURL.Scheme != "https") {
		return Config{}, errors.New("invalid ROSTERING_URL")
	}
	for _, origin := range strings.Split(value("CORS_ORIGINS", "http://localhost:5173,http://localhost:3000"), ",") {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		u, err := url.Parse(origin)
		if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") || u.Path != "" || u.RawQuery != "" || u.Fragment != "" || u.User != nil {
			return Config{}, fmt.Errorf("invalid CORS origin")
		}
		cfg.CORSOrigins = append(cfg.CORSOrigins, origin)
	}
	cfg.JWTAccessTTL, err = time.ParseDuration(value("JWT_ACCESS_TTL", "15m"))
	if err != nil || cfg.JWTAccessTTL < time.Second {
		return Config{}, errors.New("invalid JWT_ACCESS_TTL")
	}
	cfg.JWTRefreshTTL, err = time.ParseDuration(value("JWT_REFRESH_TTL", "168h"))
	if err != nil || cfg.JWTRefreshTTL <= cfg.JWTAccessTTL {
		return Config{}, errors.New("JWT_REFRESH_TTL must be greater than JWT_ACCESS_TTL")
	}
	if cfg.JWTAccessSecret != "" || cfg.JWTRefreshSecret != "" {
		if len(cfg.JWTAccessSecret) < 32 || len(cfg.JWTRefreshSecret) < 32 || cfg.JWTAccessSecret == cfg.JWTRefreshSecret {
			return Config{}, errors.New("JWT secrets must be distinct and at least 32 bytes each")
		}
	}
	cfg.SMTPPort, err = strconv.Atoi(value("SMTP_PORT", "587"))
	if err != nil || cfg.SMTPPort < 1 || cfg.SMTPPort > 65535 {
		return Config{}, errors.New("invalid SMTP_PORT")
	}
	cfg.SMTPTimeout, err = time.ParseDuration(value("SMTP_TIMEOUT", "10s"))
	if err != nil || cfg.SMTPTimeout <= 0 || cfg.SMTPTimeout > time.Minute {
		return Config{}, errors.New("SMTP_TIMEOUT must be between 0s and 1m")
	}
	if cfg.SMTPEncryption != "starttls" && cfg.SMTPEncryption != "tls" {
		return Config{}, errors.New("SMTP_ENCRYPTION must be starttls or tls")
	}
	if cfg.SMTPHost != "" || cfg.SMTPUsername != "" || cfg.SMTPPassword != "" || cfg.SMTPFrom != "" {
		if cfg.SMTPHost == "" || cfg.SMTPUsername == "" || cfg.SMTPPassword == "" || cfg.SMTPFrom == "" {
			return Config{}, errors.New("SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD and SMTP_FROM must all be configured")
		}
		if _, err := mail.ParseAddress(cfg.SMTPFrom); err != nil || strings.ContainsAny(cfg.SMTPFrom+cfg.SMTPHost, "\r\n") {
			return Config{}, errors.New("invalid SMTP_FROM or SMTP_HOST")
		}
	}
	cfg.EmailVerificationTTL, err = time.ParseDuration(value("EMAIL_VERIFICATION_TTL", "15m"))
	if err != nil || cfg.EmailVerificationTTL < time.Minute || cfg.EmailVerificationTTL > 24*time.Hour {
		return Config{}, errors.New("EMAIL_VERIFICATION_TTL must be between 1m and 24h")
	}
	cfg.EmailResendCooldown, err = time.ParseDuration(value("EMAIL_RESEND_COOLDOWN", "60s"))
	if err != nil || cfg.EmailResendCooldown < time.Second || cfg.EmailResendCooldown >= cfg.EmailVerificationTTL {
		return Config{}, errors.New("EMAIL_RESEND_COOLDOWN must be at least 1s and less than verification TTL")
	}
	verifyURL, err := url.Parse(cfg.EmailVerificationURL)
	if err != nil || verifyURL.Host == "" || (verifyURL.Scheme != "http" && verifyURL.Scheme != "https") || verifyURL.User != nil || verifyURL.RawQuery != "" || verifyURL.Fragment != "" {
		return Config{}, errors.New("invalid EMAIL_VERIFICATION_URL")
	}
	return cfg, nil
}
func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
