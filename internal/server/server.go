package server

import (
	"bufio"
	"net"

	"github.com/decoesp/tamborete/internal/config"
	"github.com/decoesp/tamborete/internal/database"
	"github.com/decoesp/tamborete/internal/persistence"
	"github.com/decoesp/tamborete/internal/resp"
)

type Server struct {
	cfg     *config.Config
	db      *database.Database
	persist *persistence.Persistence
}

func New(cfg *config.Config) *Server {
	db := database.New()
	return &Server{
		cfg:     cfg,
		db:      db,
		persist: persistence.New(cfg.Server.Persist, db),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":"+s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	parser := resp.NewParser(reader)

	for {
		_, _ = parser.Parse()
		// Processar comando
	}
}
