package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

type Config struct {
	Addr             string
	Mode             string
	DatabaseURL      string
	DBMaxConns       int32
	Location         *time.Location
	CORSOrigins      []string
	SolverURL        string
	SolverTimeout    time.Duration
	StorageDir       string
	JWTAccessSecret  string
	JWTRefreshSecret string
	JWTIssuer        string
	JWTAccessTTL     time.Duration
	JWTRefreshTTL    time.Duration
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
	return cfg, nil
}
func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
