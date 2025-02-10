package server

import (
	"Kashyap-site/server/functions"
	"log/slog"
	"net/http"
)

var (
  routes = map[string]http.HandlerFunc{
    "/": functions.Index,
  }
)

func New(port string) {
  if port == "" {
    port = "8080"
  }

  mux := http.NewServeMux()
  for route, handler := range routes {
    mux.HandleFunc(route, handler)
  }

  server := http.Server{
    Addr: ":" + port,
    Handler: mux,
  }
  
  slog.Info("Server Started on http://localhost:" + port)
  if err := server.ListenAndServe(); err != nil {
    panic(err)
  }
}
