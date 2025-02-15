package server

import (
	"net/http"
	"strings"

	"github.com/SapolovichSV/durak/auth/internal/config"
)

type Server struct {
	*http.Server
}

func New(config config.Config, mux *http.ServeMux) Server {
	server := http.Server{
		Addr:    buildServerCompleteAdrr(config.Server),
		Handler: mux,
	}
	return Server{
		Server: &server,
	}
}
func buildServerCompleteAdrr(confServer config.Server) string {
	AddrFistPart := confServer.Addr
	AddrSecondPart := confServer.Port
	res := strings.Join([]string{AddrFistPart, AddrSecondPart}, ":")
	return res
}
