package server

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkorolyov/wordofwisdom/workerpool"
	"golang.org/x/net/netutil"
)

type TCPServer struct {
	cfg Config

	handlerPool *workerpool.WorkerPool
}

type Config struct {
	Addr                string
	ConcurrentConnLimit int
	HandlerPool         workerpool.Config
}

func NewTCPServer(cfg Config) *TCPServer {
	return &TCPServer{
		cfg:         cfg,
		handlerPool: workerpool.New(cfg.HandlerPool),
	}
}

func (s *TCPServer) Serve(ctx context.Context, handler Handler) {
	listenConfig := net.ListenConfig{}

	ctx, cancel := context.WithCancel(ctx)

	l, err := listenConfig.Listen(ctx, "tcp", s.cfg.Addr)
	if err != nil {
		log.Fatalf("cant start server on addr %s: %v", s.cfg.Addr, err)
	}

	if s.cfg.ConcurrentConnLimit > 0 {
		l = netutil.LimitListener(l, s.cfg.ConcurrentConnLimit)
	}

	// pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	go s.watchShutdown(cancel, l)

	handler = connCloser(handler)

	for {
		// taking in account that epoll already used in stdlib net runtime our server already will perform well.
		// we can scale more by managing epoll explicitly e.g. https://github.com/mkorolyov/1m-go-tcp-server/blob/master/8_server_workerpool/epoll_linux.go
		conn, err := l.Accept()
		if err != nil {
			// unblocked accept after listener is closed during shutdown
			if !errors.Is(err, net.ErrClosed) {
				log.Printf("listener.Accept error: %v\n", err)
			}

			return
		}

		s.handlerPool.AddTask(func() {
			if err := handler(conn); err != nil {
				log.Printf("ERROR: handle incoming tcp conn: %v", err)
			}
		})
	}
}

func connCloser(h Handler) Handler {
	return func(conn net.Conn) error {
		defer func() {
			if err := conn.Close(); err != nil {
				log.Printf("ERROR: cant close conn: %v\n", err)
			}
		}()

		return h(conn)
	}
}

type Handler func(conn net.Conn) error

func (s *TCPServer) watchShutdown(cancelFunc context.CancelFunc, listener net.Listener) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("received %s signal from OS\n", sig.String())
	cancelFunc()

	if err := listener.Close(); err != nil {
		log.Printf("close tcp listener: %v", err)
	}

	s.handlerPool.Close()
}
