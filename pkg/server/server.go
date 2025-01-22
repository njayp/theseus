package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/njayp/theseus/pkg/manager"
)

type Server struct {
	manager *manager.Manager
}

func NewServer() (*Server, error) {
	// Initialize the image manager
	mgr, err := manager.NewManager()
	if err != nil {
		return nil, err
	}

	return &Server{
		manager: mgr,
	}, nil
}

func (s *Server) addHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the request body
	config := manager.Config{}

	// Read the body
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to add image: %s", config.ContainerConfig.Image)

	err = s.manager.AddImage(r.Context(), config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Image added and container started successfully"))
}

func (s *Server) removeHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the request body
	data := manager.RemoveRequest{}

	// Read the body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received request to remove image: %s", data.ImageName)

	err = s.manager.RemoveImage(r.Context(), data.ImageName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Image removed successfully"))
}

func (s *Server) upgradeHandler(w http.ResponseWriter, r *http.Request) {
	// Define a struct to hold the request body
	data := BuildPayload{}

	// Read the body
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the image name from the struct
	imageName := data.Repository.RepoName
	log.Printf("Received request to upgrade image: %s", imageName)

	err = s.manager.UpgradeImage(r.Context(), imageName)
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
