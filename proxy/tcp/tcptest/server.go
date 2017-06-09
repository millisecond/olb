package tcptest

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/fabiolb/fabio/proxy/internal"
	"github.com/fabiolb/fabio/proxy/tcp"
)

// Server is a TCP test server that binds to a random port.
type Server struct {
	// Addr is the address the server is listening on in the form ipaddr:port.
	Addr     string
	Listener net.Listener

	// TLS is the optional TLS configuration, populated with a new config
	// after TLS is started. If set on an unstarted server before StartTLS
	// is called, existing fields are copied into the new config.
	TLS *tls.Config

	// Config may be changed after calling NewUnstartedServer and
	// before Start or StartTLS.
	Config *tcp.Server

	// srv is the actual running server.
	srv *tcp.Server
}

func (s *Server) Start() {
	if s.Addr != "" {
		panic("Server already started")
	}

	s.Addr = s.Listener.Addr().String()
	s.srv = new(tcp.Server)
	*s.srv = *s.Config
	s.srv.Addr = s.Addr
	go s.srv.Serve(s.Listener)
}

func (s *Server) StartTLS() {
	if s.Addr != "" {
		panic("Server already started")
	}

	s.Addr = s.Listener.Addr().String()
	s.srv = new(tcp.Server)
	*s.srv = *s.Config
	s.srv.Addr = s.Addr

	cert, err := tls.X509KeyPair(internal.LocalhostCert, internal.LocalhostKey)
	if err != nil {
		panic(fmt.Sprintf("tcptest: NewTLSServer: %v", err))
	}

	existingConfig := s.TLS
	if existingConfig != nil {
		s.TLS = existingConfig.Clone()
	} else {
		s.TLS = new(tls.Config)
	}
	if len(s.TLS.Certificates) == 0 {
		s.TLS.Certificates = []tls.Certificate{cert}
	}
	s.Listener = tls.NewListener(s.Listener, s.TLS)
	go s.srv.Serve(s.Listener)
}

func (s *Server) Close() error {
	if s.Addr == "" {
		panic("Server not started")
	}
	return s.srv.Close()
}

func NewServer(h tcp.Handler) *Server {
	srv := NewUnstartedServer(h)
	srv.Start()
	return srv
}

func NewTLSServer(h tcp.Handler) *Server {
	srv := NewUnstartedServer(h)
	srv.StartTLS()
	return srv
}

func NewUnstartedServer(h tcp.Handler) *Server {
	return &Server{
		Listener: newLocalListener(),
		Config:   &tcp.Server{Handler: h},
	}
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		l, err = net.Listen("tcp6", "[::1]:0")
		if err != nil {
			panic("tcptest: Failed to listen on a port: " + err.Error())
		}
	}
	return l
}
