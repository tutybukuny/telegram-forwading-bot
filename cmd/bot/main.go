package main

import (
	"forwarding-bot/config"
	"forwarding-bot/pkg/container"
	handleossignal "forwarding-bot/pkg/handle-os-signal"
	"forwarding-bot/pkg/l"
	"forwarding-bot/pkg/l/sentry"
)

func main() {
	ll := l.New()
	cfg := config.Load(ll)

	if cfg.SentryConfig.Enabled {
		ll = l.NewWithSentry(&sentry.Configuration{
			DSN: cfg.SentryConfig.DNS,
			Trace: struct{ Disabled bool }{
				Disabled: !cfg.SentryConfig.Trace,
			},
		})
	}

	container.NamedSingleton("ll", func() l.Logger {
		return ll
	})

	// init os signal handle
	shutdown := handleossignal.New(ll)
	shutdown.HandleDefer(func() {
		ll.Sync()
	})
	container.NamedSingleton("shutdown", func() handleossignal.IShutdownHandler {
		return shutdown
	})

	bootstrap(cfg)

	go startBot(cfg)

	// handle signal
	if cfg.Environment == "D" {
		shutdown.SetTimeout(1)
	}
	shutdown.Handle()
}
