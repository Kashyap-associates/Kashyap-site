package main

import (
	"Kashyap-site/config"
	"Kashyap-site/server"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	server_port, admin_port, tg_token := config.New()
	go func(admin_port string) {
		server_admin := http.Server{
			Addr:    ":" + admin_port,
			Handler: server.NewAdmin(),
		}
		slog.Info("Admin Server Started", "port", server_admin.Addr, "url", "http://localhost"+server_admin.Addr)
		if err := server_admin.ListenAndServe(); err != nil {
			slog.Error("Server Error", "msg:", err)
			os.Exit(1)
		}
	}(admin_port)

	go func(token string) {
		server.Telegram(token)
	}(tg_token)

	serv := http.Server{
		Addr:    ":" + server_port,
		Handler: server.New(),
	}

	slog.Info("Server Started", "port", serv.Addr, "url", "http://localhost"+serv.Addr)
	if err := serv.ListenAndServe(); err != nil {
		slog.Error("Server Error", "msg:", err)
		os.Exit(1)
	}
}
