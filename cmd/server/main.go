package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/decoesp/tamborete/internal/auth"
	"github.com/decoesp/tamborete/internal/config"
	"github.com/decoesp/tamborete/internal/database"
	"github.com/decoesp/tamborete/internal/persistence"
	"github.com/decoesp/tamborete/internal/resp"
)

type Server struct {
	cfg     *config.Config
	db      *database.Database
	auth    *auth.Auth
	persist *persistence.Persistence
}

func New(cfg *config.Config) *Server {
	db := database.New()
	return &Server{
		cfg:     cfg,
		db:      db,
		auth:    auth.New(cfg.Server.Auth),
		persist: persistence.New(cfg.Server.Persist, db),
	}
}

func (s *Server) Start() error {
	if err := s.persist.Load(); err != nil {
		return fmt.Errorf("failed to load data: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	go s.handleSignals()
	fmt.Printf("Server listening on port %s\n", s.cfg.Server.Port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	authenticated := false
	parser := resp.NewParser(bufio.NewReader(conn))

	for {
		command, err := parser.Parse()
		if err != nil {
			return
		}

		args, ok := command.([]interface{})
		if !ok || len(args) == 0 {
			conn.Write(resp.Serialize("ERR invalid command"))
			continue
		}

		cmd := strings.ToUpper(args[0].(string))
		params := args[1:]

		if !authenticated && cmd != "AUTH" {
			conn.Write(resp.Serialize("NOAUTH Authentication required"))
			continue
		}

		switch cmd {
		case "AUTH":
			if len(params) != 1 {
				conn.Write(resp.Serialize("ERR wrong number of arguments"))
				continue
			}
			if err := s.auth.Authenticate(params[0].(string)); err != nil {
				conn.Write(resp.Serialize(err))
			} else {
				authenticated = true
				conn.Write(resp.Serialize("OK"))
			}

		case "SET":
			if len(params) != 2 {
				conn.Write(resp.Serialize("ERR wrong number of arguments"))
				continue
			}
			s.db.Set(params[0].(string), params[1].(string))
			conn.Write(resp.Serialize("OK"))

		case "GET":
			if len(params) != 1 {
				conn.Write(resp.Serialize("ERR wrong number of arguments"))
				continue
			}
			val, exists := s.db.Get(params[0].(string))
			if !exists {
				conn.Write(resp.Serialize(nil))
				continue
			}
			conn.Write(resp.Serialize(val))

		case "LPUSH":
			if len(params) < 2 {
				conn.Write(resp.Serialize("ERR wrong number of arguments"))
				continue
			}
			count := 0
			for _, value := range params[1:] {
				count = s.db.LPush(params[0].(string), value.(string))
			}
			conn.Write(resp.Serialize(count))

		case "SAVE":
			if err := s.persist.Save(); err != nil {
				conn.Write(resp.Serialize(err))
			} else {
				conn.Write(resp.Serialize("OK"))
			}

		default:
			conn.Write(resp.Serialize("ERR unknown command"))
		}
	}
}

func (s *Server) handleSignals() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nSaving data before shutdown...")
	if err := s.persist.Save(); err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	}
	os.Exit(0)
}

func main() {
	// TODO: Add proper server initialization
	fmt.Println("Server starting...")
}
