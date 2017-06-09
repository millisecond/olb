package admin

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/fabiolb/fabio/admin/api"
	"github.com/fabiolb/fabio/admin/ui"
	"github.com/fabiolb/fabio/config"
	"github.com/fabiolb/fabio/proxy"
)

// Server provides the HTTP server for the admin UI and API.
type Server struct {
	Access   string
	Color    string
	Title    string
	Version  string
	Commands string
	Cfg      *config.Config
}

// ListenAndServe starts the admin server.
func (s *Server) ListenAndServe(l config.Listen, tlscfg *tls.Config) error {
	return proxy.ListenAndServeHTTP(l, s.handler(), tlscfg)
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()

	switch s.Access {
	case "ro":
		mux.HandleFunc("/api/manual", forbidden)
		mux.HandleFunc("/manual", forbidden)
	case "rw":
		mux.Handle("/api/manual", &api.ManualHandler{})
		mux.Handle("/manual", &ui.ManualHandler{Color: s.Color, Title: s.Title, Version: s.Version, Commands: s.Commands})
	}

	mux.Handle("/api/config", &api.ConfigHandler{s.Cfg})
	mux.Handle("/api/routes", &api.RoutesHandler{})
	mux.Handle("/api/version", &api.VersionHandler{s.Version})
	mux.Handle("/routes", &ui.RoutesHandler{Color: s.Color, Title: s.Title, Version: s.Version})
	mux.HandleFunc("/logo.svg", ui.HandleLogo)
	mux.HandleFunc("/health", handleHealth)
	mux.Handle("/", http.RedirectHandler("/routes", http.StatusSeeOther))
	return mux
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func forbidden(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Forbidden", http.StatusForbidden)
}
