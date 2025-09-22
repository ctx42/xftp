package ftpsrv

import (
	"crypto/tls"
	"net"
)

// Server represents FTP server.
type Server struct {
	cfg      Config            // FTP server configuration.
	certs    []tls.Certificate // TLS certificates.
	listener net.Listener      // FTP server listener.
}

// NewServer returns a new instance of [Server].
func NewServer(cfg Config) (*Server, error) {
	srv := &Server{cfg: cfg}
	if srv.cfg.cert != nil && srv.cfg.key != nil {
		var err error
		srv.certs = make([]tls.Certificate, 1)
		srv.certs[0], err = tls.X509KeyPair(srv.cfg.cert, srv.cfg.key)
		if err != nil {
			return nil, err
		}
	}
	return srv, nil
}

// ListenAndServe is a blocking call that starts the FTP server.
func (srv Server) ListenAndServe() error {
	return nil
}

// Close closes the FTP server.
func (srv Server) Close() error {
	return nil
}
