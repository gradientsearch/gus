package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/gradientsearch/gus/api/cmd/services/gus/build/all"
	"github.com/gradientsearch/gus/api/http/api/debug"
	"github.com/gradientsearch/gus/api/http/api/mux"
	"github.com/gradientsearch/gus/app/api/authclient"
	"github.com/gradientsearch/gus/business/api/sqldb"
	"github.com/gradientsearch/gus/business/domain/chatbus"
	"github.com/gradientsearch/gus/business/domain/chatbus/llms"
	"github.com/gradientsearch/gus/business/domain/chatbus/llms/llama"
	"github.com/gradientsearch/gus/business/domain/chatbus/stores/chatdb"
	"github.com/gradientsearch/gus/business/domain/userbus"
	"github.com/gradientsearch/gus/business/domain/userbus/stores/userdb"
	"github.com/gradientsearch/gus/foundation/logger"
)

var build = "develop"

func main() {
	l := logger.New(os.Stdout, logger.LevelInfo, "GUS", nil)
	if err := run(context.Background(), l); err != nil {
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	cfg := struct {
		conf.Version
		Web struct {
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			APIHost            string        `conf:"default:0.0.0.0:3000"`
			DebugHost          string        `conf:"default:0.0.0.0:3010"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
		Auth struct {
			Host string `conf:"default:http://auth-service.gus-system.svc.cluster.local:6000"`
		}
		DB struct {
			User         string `conf:"default:postgres"`
			Password     string `conf:"default:postgres,mask"`
			HostPort     string `conf:"default:database-service.gus-system.svc.cluster.local"`
			Name         string `conf:"default:postgres"`
			MaxIdleConns int    `conf:"default:2"`
			MaxOpenConns int    `conf:"default:0"`
			DisableTLS   bool   `conf:"default:true"`
		}
		LLM struct {
			LLMHost string `conf:"default:http://localllm.dev:11434"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "GUS",
		},
	}

	const prefix = "GUS"
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// ---------------------------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Info(ctx, "startup", "config", out)

	expvar.NewString("build").Set(cfg.Build)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "hostport", cfg.DB.HostPort)

	db, err := sqldb.Open(sqldb.Config{
		User:         cfg.DB.User,
		Password:     cfg.DB.Password,
		HostPort:     cfg.DB.HostPort,
		Name:         cfg.DB.Name,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		DisableTLS:   cfg.DB.DisableTLS,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	// -------------------------------------------------------------------------
	// Initialize authentication support

	log.Info(ctx, "startup", "status", "initializing authentication support")

	logFunc := func(ctx context.Context, msg string, v ...any) {
		log.Info(ctx, msg, v...)
	}
	authClient := authclient.New(cfg.Auth.Host, logFunc)
	// -------------------------------------------------------------------------
	// Create LLM

	_ = &llama.Llama{
		BaseURL: "http://localllm.dev:11434",
		Client:  &http.Client{},
		Model:   "llama3.2",
		Stream:  false,
	}

	mockLlm := &llms.Mock{}
	// -------------------------------------------------------------------------
	// Create Business Packages

	userBus := userbus.NewBusiness(log, userdb.NewStore(log, db))
	chatBus := chatbus.NewBusiness(log, chatdb.NewStore(log, db), mockLlm)

	defer db.Close()
	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "msg", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGINT)

	cfgMux := mux.Config{
		Build:      build,
		AuthClient: authClient,
		Log:        log,
		ChatBus:    chatBus,
		UserBus:    userBus,
	}

	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      mux.WebAPI(cfgMux, all.Routes()),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     logger.NewStdLogger(log, logger.LevelError),
	}

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", api.Addr)

		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
