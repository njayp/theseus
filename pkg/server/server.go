package server

import (
	"fmt"
	"net/http"

	"github.com/njayp/theseus/pkg/manager"
)

type Server struct {
	manager *manager.ImageManager
}

func NewServer() (*Server, error) {
	// Initialize the image manager
	mgr, err := manager.NewImageManager()
	if err != nil {
		return nil, err
	}

	return &Server{
		manager: mgr,
	}, nil
}

func (s *Server) addHandler(w http.ResponseWriter, r *http.Request) {
	err := s.manager.AddImage(r.Context(), "nginx:latest", "nginx-container")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Image added and container started successfully"))
}

func (s *Server) removeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.manager.RemoveImage(r.Context(), "nginx:latest")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Image removed successfully"))
}

func (s *Server) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	err := s.manager.UpgradeImage(r.Context(), "nginx:latest")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Image upgraded successfully"))
}

func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

func (s *Server) Start(port int) error {
	http.HandleFunc("/add", s.addHandler)
	http.HandleFunc("/remove", s.removeHandler)
	http.HandleFunc("/upgrade", s.upgradeHandler)
	http.HandleFunc("/", s.notFoundHandler) // Catch-all for 404
	addr := fmt.Sprintf(":%d", port)
	return http.ListenAndServe(addr, nil)
}
