package config

import "testing"

func clearConfig(t *testing.T) {
	for _, key := range []string{"SERVER_ADDR", "GIN_MODE", "DATABASE_URL", "DB_MAX_CONNS", "APP_TIMEZONE", "CORS_ORIGINS", "ROSTERING_URL", "ROSTERING_TIMEOUT", "STORAGE_DIR", "JWT_ACCESS_SECRET", "JWT_REFRESH_SECRET", "JWT_ISSUER", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL", "SMTP_HOST", "SMTP_PORT", "SMTP_USERNAME", "SMTP_PASSWORD", "SMTP_FROM", "SMTP_ENCRYPTION", "SMTP_TIMEOUT", "EMAIL_VERIFICATION_URL", "EMAIL_VERIFICATION_TTL", "EMAIL_RESEND_COOLDOWN"} {
		t.Setenv(key, "")
	}
}
func TestConfigDefaults(t *testing.T) {
	clearConfig(t)
	cfg, err := FromEnv()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Addr != "127.0.0.1:8081" || cfg.Location.String() != "Asia/Ho_Chi_Minh" || cfg.DBMaxConns != 10 {
		t.Fatalf("%+v", cfg)
	}
}
func TestInvalidConfig(t *testing.T) {
	for _, tc := range []struct{ key, value string }{
		{"SERVER_ADDR", ":0"}, {"SERVER_ADDR", "bad"}, {"GIN_MODE", "prod"},
		{"DB_MAX_CONNS", "0"}, {"APP_TIMEZONE", "wrong/zone"},
		{"ROSTERING_URL", "file:///tmp"}, {"ROSTERING_TIMEOUT", "0s"},
		{"CORS_ORIGINS", "*"}, {"CORS_ORIGINS", "http://localhost:3000/path"},
	} {
		t.Run(tc.key+tc.value, func(t *testing.T) {
			clearConfig(t)
			t.Setenv(tc.key, tc.value)
			if _, err := FromEnv(); err == nil {
				t.Fatal("invalid config accepted")
			}
		})
	}
}
