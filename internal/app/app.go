package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"abclinic.local/backend/internal/attachment"
	"abclinic.local/backend/internal/auth"
	"abclinic.local/backend/internal/booking"
	"abclinic.local/backend/internal/closure"
	"abclinic.local/backend/internal/config"
	"abclinic.local/backend/internal/doctor"
	"abclinic.local/backend/internal/examination"
	"abclinic.local/backend/internal/jobs"
	"abclinic.local/backend/internal/leave"
	"abclinic.local/backend/internal/medicine"
	"abclinic.local/backend/internal/notification"
	"abclinic.local/backend/internal/patient"
	"abclinic.local/backend/internal/platform/clock"
	"abclinic.local/backend/internal/platform/middleware"
	"abclinic.local/backend/internal/platform/postgres"
	"abclinic.local/backend/internal/platform/swagger"
	"abclinic.local/backend/internal/roster"
	"abclinic.local/backend/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Config config.Config
	DB     *pgxpool.Pool
	Clock  clock.Clock
	Auth   *auth.Service
}

func New(ctx context.Context, cfg config.Config) (*App, error) {
	a := &App{Config: cfg, Clock: clock.System{Location: cfg.Location}}
	if cfg.DatabaseURL != "" {
		pool, err := postgres.Open(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
		if err != nil {
			return nil, err
		}
		a.DB = pool
	} else {
		slog.Warn("DATABASE_URL is empty; API docs work but readiness returns 503")
	}
	var tokens *auth.Tokens
	if cfg.JWTAccessSecret != "" {
		var err error
		tokens, err = auth.NewTokens(auth.TokenConfig{
			AccessSecret: cfg.JWTAccessSecret, RefreshSecret: cfg.JWTRefreshSecret, Issuer: cfg.JWTIssuer,
			AccessTTL: cfg.JWTAccessTTL, RefreshTTL: cfg.JWTRefreshTTL,
		}, a.Clock)
		if err != nil {
			a.Close()
			return nil, err
		}
	} else {
		slog.Warn("JWT secrets are missing; run scripts/setup.ps1 to configure authentication")
	}
	a.Auth = auth.NewService(auth.NewPostgresRepository(a.DB), tokens)
	return a, nil
}
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
func (a *App) Router() *gin.Engine {
	gin.SetMode(a.Config.Mode)
	r := gin.New()
	r.Use(middleware.RequestID(), gin.Logger(), gin.Recovery(), middleware.CORS(a.Config.CORSOrigins))
	_ = r.SetTrustedProxies(nil)
	swagger.Register(r)
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok", "mode": "scaffold"}) })
	r.GET("/readyz", func(c *gin.Context) {
		if a.DB == nil {
			c.JSON(503, gin.H{"status": "not_ready", "reason": "database_not_configured"})
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := a.DB.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"status": "not_ready", "reason": "database_unavailable"})
			return
		}
		c.JSON(200, gin.H{"status": "ready", "businessModules": "auth_implemented"})
	})
	api := r.Group("/api")
	auth.NewHandler(a.Auth).RegisterRoutes(api)
	user.NewHandler(user.NewService(a.Clock)).RegisterRoutes(api)
	doctor.NewHandler(doctor.NewService(a.Clock)).RegisterRoutes(api)
	patient.NewHandler(patient.NewService(a.Clock)).RegisterRoutes(api)
	medicine.NewHandler(medicine.NewService(a.Clock)).RegisterRoutes(api)
	examination.NewHandler(examination.NewService(a.Clock)).RegisterRoutes(api)
	roster.NewHandler(roster.NewService(a.Clock)).RegisterRoutes(api)
	leave.NewHandler(leave.NewService(a.Clock)).RegisterRoutes(api)
	closure.NewHandler(closure.NewService(a.Clock)).RegisterRoutes(api)
	booking.NewHandler(booking.NewService(a.Clock)).RegisterRoutes(api)
	notification.NewHandler(notification.NewService(a.Clock)).RegisterRoutes(api)
	attachment.NewHandler(attachment.NewService(a.Clock)).RegisterRoutes(api)
	r.NoRoute(func(c *gin.Context) { c.JSON(404, gin.H{"message": "Route not found"}) })
	return r
}
func (a *App) RunAPI(ctx context.Context) error {
	srv := &http.Server{Addr: a.Config.Addr, Handler: a.Router(),
		ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	errCh := make(chan error, 1)
	go func() { errCh <- srv.ListenAndServe() }()
	slog.Info("ABClinic API starting", "addr", a.Config.Addr, "swagger", "/swagger/index.html")
	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		shutdown, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			_ = srv.Close()
			return err
		}
		err := <-errCh
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
func (a *App) RunWorker(ctx context.Context) error { return jobs.Run(ctx) }
