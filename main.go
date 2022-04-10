package main

import (
	"arion_shot_api/app"
	"arion_shot_api/platform/conf"
	"arion_shot_api/utils/cloudinary"
	"context"
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arion_shot_api/platform/auth"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

func main() {
	if err := run(); err != nil {
		log.Default().Fatal(err)
	}
}

func run() error {
	logger := log.New(os.Stdout, "ArionShots : ", log.LstdFlags)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8001" // Default port if not specified
	}

	var cfg struct {
		Web struct {
			Address         string        `conf:"default:localhost:8001"`
			Debug           string        `conf:"default:localhost:6060"`
			ReadTimeout     time.Duration `conf:"default:5s"`
			WriteTimeout    time.Duration `conf:"default:5s"`
			ShutdownTimeout time.Duration `conf:"default:5s"`
		}
		Auth struct {
			KeyID          string `conf:"default:1"`
			PrivateKeyFile string `conf:"default:private.pem"`
			Algorithm      string `conf:"default:RS256"`
		}
	}

	if err := conf.Parse(os.Args[1:], "ArionShots", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("ArionShots", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting
	logger.Printf("main: Started")
	defer logger.Println("main: Completed")

	cloudinary.Initialize()

	// =========================================================================
	// Initialize authentication support
	authenticator, err := createAuth(cfg.Auth.PrivateKeyFile, cfg.Auth.KeyID, cfg.Auth.Algorithm)
	if err != nil {
		return err
	}

	// =========================================================================
	// Start API Service
	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         "0.0.0.0:8080",
		Handler:      app.API(shutdown, logger, authenticator),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverError := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		logger.Printf("main : API listening on %s", api.Addr)
		serverError <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverError:
		return errors.Wrap(err, "Listening and serving")

	case sig := <-shutdown:
		logger.Println("main: Start shutdown", sig)

		// Give outstanding requests a deadline for completion.
		timeout := cfg.Web.ShutdownTimeout
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		defer cancel()

		// Asking listener to shut down and load shed.
		err := api.Shutdown(ctx)

		if err != nil {
			logger.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			err = api.Close()
		}

		if err != nil {
			return errors.Wrap(err, "graceful shutdown")
		}

		if sig == syscall.SIGSTOP {
			return errors.New("Integrity error detected, asking for self shutdown")
		}
	}

	return nil
}

func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {
	keyContents, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading auth private key")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing auth private key")
	}

	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	return auth.NewAuthenticator(key, keyID, algorithm, public)
}
