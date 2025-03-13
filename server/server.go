package server

import (
	"Kashyap-site/server/functions"
	"Kashyap-site/server/middleware"
	"log/slog"
	"net/http"
	"os"
)

var routes = map[string]http.HandlerFunc{
	"/":       functions.Index,
	"/admin":  functions.Admin,
	// "/thanks": functions.Thanks,
	// "/404":    functions.NotFound,
	// "/error":  functions.Error,
}

func New(port string) {
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	for route, handler := range routes {
		mux.HandleFunc(route, handler)
	}

	handler := middleware.LoggingMiddleware(mux)
	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	slog.Info("Server Started on http://localhost:" + port)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("Server Error", "msg:", err)
		os.Exit(1)
	}
}
