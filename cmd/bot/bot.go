package main

import (
	"forwarding-bot/config"
	"forwarding-bot/internal/listener"
	"forwarding-bot/pkg/container"
	handleossignal "forwarding-bot/pkg/handle-os-signal"
	"forwarding-bot/pkg/l"
)

func startBot(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	teleListener := listener.New(cfg)

	ll.Info("start bot")
	teleListener.Listen()
}
