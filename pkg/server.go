package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nate-trojian/ccg-game-server/internal"
	"github.com/nate-trojian/ccg-game-server/pkg/matchmaking"
	"go.uber.org/zap"
)

// Server - Game Server instance
type Server struct {
	logger *zap.Logger
	upgrader *websocket.Upgrader
	server *http.Server
	registerClient chan *Client
}

// NewServer - Create a new Server
func NewServer(rc chan *Client) *Server {
	server := &Server{
		logger: internal.NewLogger("server"),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		registerClient: rc,
	}

	// TODO - This should take a config object
	mux := http.NewServeMux()
	mux.HandleFunc("/login", server.login)
	mux.HandleFunc("/match", server.requestMatch)
	mux.HandleFunc("/ws", server.ws)
	s := &http.Server{
		Addr: ":8080",
		Handler: mux,
	}

	server.server = s

	return server
}

// Start - Start the Game Server
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}
	return nil
}

// Shutdown - Shutdown the Game Server
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *Server) requestMatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusText(http.StatusMethodNotAllowed), errors.New("bad method"), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.writeError(w, "Bad Matchmaking Request Body", err, http.StatusBadRequest)
		return
	}

	var request matchmaking.Request
	err = json.Unmarshal(body, &request)
	if err != nil {
		s.writeError(w, "Bad Matchmaking Request", err, http.StatusBadRequest)
		return
	}

	// TODO - Add to matckmaking

	w.WriteHeader(http.StatusOK)
}

func (s *Server) writeError(w http.ResponseWriter, msg string, err error, header int) {
	s.logger.Error(msg, zap.Error(err))
	w.WriteHeader(header)
	resp := ErrorResponse{
		Error: err.Error(),
		Message: msg,
	}
	body, _ := json.Marshal(&resp)
	_, _  = w.Write(body)
}

func (s *Server) ws(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.registerClient <- NewClient(r.RemoteAddr, conn)
}