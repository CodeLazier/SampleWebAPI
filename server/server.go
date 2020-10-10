package server

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

//wrap http.Server
type Server struct {
	*http.Server
	// Defaults to 0, indicating no limit.
	MaxConnections int
}

func (srv *Server) Listen() (net.Listener, error) {
	//for unix/linux

	//var ln net.Listener
	//if strings.HasPrefix(srv.Addr, "systemd:") {
	//	name := srv.Addr[8:]
	//	listeners, _ := activation.ListenersWithNames()
	//	listener, ok := listeners[name]
	//	if !ok {
	//		return nil, fmt.Errorf("listen systemd %s: socket not found", name)
	//	}
	//	ln = listener[0]
	//} else {
	//var err error
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		return nil, err
	}
	//}
	//if srv.MaxConnections > 0 {
	//	ln = netutil.LimitListener(ln, srv.MaxConnections)
	//}

	return ln, nil
}

func NewTLS(addr string, cert tls.Certificate, handler http.Handler) *Server {
	if addr == "" {
		addr = ":https"
	}
	srv := New(addr, handler)
	srv.TLSConfig.Certificates = []tls.Certificate{cert}

	return srv
}

func (srv *Server) listenAndServe() error {
	ln, err := srv.Listen()
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

func (srv *Server) listenAndServeTLS(certFile, keyFile string) error {
	ln, err := srv.Listen()
	if err != nil {
		return err
	}
	return srv.ServeTLS(ln, certFile, keyFile)
}

func (srv *Server) Start() error {
	ln, err := srv.Listen()
	if err != nil {
		return err
	}
	if srv.isTLS() {
		ln = tls.NewListener(ln, srv.TLSConfig)
	}
	return srv.Serve(ln)
}

func (srv *Server) isTLS() bool {
	return len(srv.TLSConfig.Certificates) > 0 || srv.TLSConfig.GetCertificate != nil
}

func New(addr string, handler http.Handler) *Server {
	if addr == "" {
		addr = ":http"
	}
	srv := &Server{
		Server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
			IdleTimeout:  120 * time.Second,
			TLSConfig: &tls.Config{
				NextProtos:       []string{"h2", "http/1.1"},
				MinVersion:       tls.VersionTLS12,
				CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
				CipherSuites: []uint16{
					tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
					tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
					tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				},
				PreferServerCipherSuites: true,
			},
		},
	}

	return srv
}
