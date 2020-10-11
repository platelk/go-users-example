package main

import (
	"go-users-example/domain/users"
	"go-users-example/infra/logger"
	"go-users-example/infra/pwdhasher"
	"go-users-example/infra/usernotifier"
	"go-users-example/infra/userstore"
	"go-users-example/transport/http"
)

func main() {
	// Load configuration from different sources
	cfg := Load()

	// Initialise log
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		panic(err)
	}
	log.Debug().Interface("config", cfg).Send()

	// Initialise user store
	usrStore := userstore.NewInMemory()

	// Initialise user notifier
	usrNotifier := usernotifier.NewInMemory()
	go func(c chan *users.ChangeEvent) {
		for e := range c {
			log.Info().Interface("change-event", e).Msg("receive event on user change")
		}
	}(usrNotifier.Listen())

	// Build http server
	srv := http.NewBuilder(log, cfg.HTTP).
		WithV1CreateUser(users.SetupCreate(log, usrNotifier, usrStore, pwdhasher.NewBcrypt())).
		WithV1UpdateUser(users.SetupUpdate(log, usrNotifier, usrStore)).
		WithV1DeleteUser(users.SetupDelete(log, usrNotifier, usrStore)).
		WithV1SearchUser(users.SetupSearch(log, usrStore)).
		WithHealthCheck().
		Build()

	// Run HTTP server
	log.Info().Msg("running http server")
	if err := srv.Run(); err != nil {
		log.Fatal().Err(err).Msg("http server didn't end correctly")
	}
	log.Info().Msg("done")
}
