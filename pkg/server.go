package pkg

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
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
	// TODO - This should take a config object
	server := &Server{
		logger: internal.NewLogger("server"),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		registerClient: rc,
	}

	r := mux.NewRouter()
	// Auth Service
	// TODO - Add
	// User Service
	// TODO - Add
	// Matchmaking Service
	r.HandleFunc("/match", server.requestMatch).Methods(http.MethodPost)
	// Game Service
	r.HandleFunc("/{gameId}/join", server.join).Methods(http.MethodGet).Queries("player_id", "")
	r.HandleFunc("/{gameId}/ws", server.ws).Methods(http.MethodGet)
	s := &http.Server{
		Addr: ":8080",
		Handler: r,
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

	// TODO - Add to matchmaking

	w.WriteHeader(http.StatusOK)
}

func (s *Server) join(w http.ResponseWriter, r *http.Request) {
	// TODO - Validate user joining
	http.Redirect(w, r, "/ws", http.StatusTemporaryRedirect)
}

func (s *Server) ws(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	s.registerClient <- NewClient(r.RemoteAddr, conn)
}