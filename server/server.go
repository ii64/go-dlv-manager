package server

import (
	"context"
	"net"
	"time"

	"github.com/ii64/go-dlv-manager/log"
	"github.com/rs/zerolog"
)

type Server struct {
	lis    net.Listener
	logger zerolog.Logger
	opt    *Option
}

func New(opt *Option) (*Server, error) {
	err := opt.Validate()
	if err != nil {
		return nil, err
	}

	s := &Server{}
	s.opt = opt
	s.logger = log.Logger()
	return s, nil
}

func (s *Server) Listen(addr string) (err error) {
	s.lis, err = net.Listen("tcp", addr)
	return
}

func (s *Server) Serve(ctx context.Context) (err error) {
	var conn net.Conn
	for {
		conn, err = s.lis.Accept()
		if err != nil {
			s.logger.Err(err).Msg("accept error")
			break
		}

		// conn = s.connMod(conn)
		go s.handleConnection(ctx, conn)
	}
	return
}

func (s *Server) connMod(conn net.Conn) net.Conn {
	conn.SetDeadline(time.Now().Add(time.Second * 20))
	return conn
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	var err error
	defer conn.Close()

	// create new debugger IPC.
	ipc, err := newIPC(s.opt.Program, s.opt.ProgramArgs)
	if err != nil {
		s.logger.Err(err).Msg("handle connection ipc")
		return
	}
	defer ipc.Close()
	// dial to debugger tcp.
	remote, err := ipc.Dial()
	if err != nil {
		s.logger.Err(err).Msg("handle connection ipc dial")
		return
	}
	defer remote.Close()

	// client <- IPC server
	go func() { // late cleanup
		err := s.createStream(remote, conn, false)
		if err != nil {
			s.logger.Err(err).Msg("stream client <- IPC server")
		}
	}()

	// client -> IPC server
	err = s.createStream(conn, remote, true)
	if err != nil {
		s.logger.Err(err).Msg("stream client -> IPC server")
	}
	return
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) (err error) {
	s.logger.Info().Str("addr", addr).Msg("serving")
	err = s.Listen(addr)
	if err != nil {
		return
	}
	return s.Serve(ctx)
}

func (s *Server) Close() error {
	if s.lis != nil {
		return s.lis.Close()
	}
	return nil
}
