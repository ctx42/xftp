package ftpsrv

import (
	"context"
	"crypto/tls"
	"net"
	"strconv"

	"github.com/rs/zerolog"
)

// Server represents FTP server.
type Server struct {
	cfg    Config       // FTP server configuration.
	tlsCfg *tls.Config  // TLS configuration.
	lnr    net.Listener // FTP server lnr.
	cxl    func()
	log    zerolog.Logger
}

// NewServer returns a new instance of [Server].
func NewServer(cfg Config) (*Server, error) {
	srv := &Server{
		cfg: cfg,
	}
	if srv.cfg.cert != nil && srv.cfg.key != nil {
		cert, err := tls.X509KeyPair(srv.cfg.cert, srv.cfg.key)
		if err != nil {
			return nil, err
		}
		srv.tlsCfg = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			ClientAuth:         tls.NoClientCert,
			InsecureSkipVerify: true, // nolint: gosec
		}
	}
	return srv, nil
}

// ListenAndServe is a blocking call that starts the FTP server.
func (srv *Server) ListenAndServe(ctx context.Context) error {
	address := net.JoinHostPort(srv.cfg.host, strconv.Itoa(srv.cfg.port))

	var err error
	if srv.tlsCfg != nil && srv.cfg.implicitTLS {
		srv.lnr, err = tls.Listen("tcp", address, srv.tlsCfg)
	} else {
		srv.lnr, err = net.Listen("tcp", address)
	}
	if err != nil {
		return err
	}
	return srv.Serve(ctx, srv.lnr)
}

func (srv *Server) Serve(ctx context.Context, lnr net.Listener) error {
	srv.lnr = lnr
	for {
		conn, err := srv.lnr.Accept()
		if err = srv.handleAcceptError(err); err != nil {
			return err
		}
		ses := NewSession(srv.cfg, conn.LocalAddr(), conn.RemoteAddr())
		cc := NewCtrlCon(ses, conn, srv.log).WithTLS(srv.tlsCfg)
		go cc.Listen()
	}
}

func (srv *Server) handleAcceptError(err error) error {
	return err
}

// Close closes the FTP server.
func (srv *Server) Close() error {
	return nil
}
