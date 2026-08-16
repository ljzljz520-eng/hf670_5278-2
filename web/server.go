package web

import (
	"embed"
	"encoding/json"
	"net/http"
	"strings"

	"contractseal/internal/contract"
	"github.com/go-chi/chi/v5"
)

//go:embed page.html
var pageFS embed.FS

type Server struct {
	service *contract.Service
}

func NewHandler(service *contract.Service) http.Handler {
	server := &Server{service: service}
	router := chi.NewRouter()
	router.Get("/", server.page)
	router.Post("/api/contract-seals", server.submit)
	router.Post("/api/contract-seals/{id}/process", server.process)
	router.Get("/api/contract-seals/{id}", server.get)
	return router
}

func (s *Server) page(w http.ResponseWriter, r *http.Request) {
	data, err := pageFS.ReadFile("page.html")
	if err != nil {
		http.Error(w, "page unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *Server) submit(w http.ResponseWriter, r *http.Request) {
	var input contract.Submission
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.service.Submit(input)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

type processRequest struct {
	Files []contract.FileInput `json:"files"`
}

func (s *Server) process(w http.ResponseWriter, r *http.Request) {
	var input processRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.service.Process(chi.URLParam(r, "id"), input.Files)
	if err != nil {
		status := http.StatusBadRequest
		if strings.HasSuffix(err.Error(), "not found") {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (s *Server) get(w http.ResponseWriter, r *http.Request) {
	item, ok := s.service.Get(chi.URLParam(r, "id"))
	if !ok {
		writeError(w, http.StatusNotFound, nil)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func decodeJSON(r *http.Request, target any) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, err error) {
	message := "not found"
	if err != nil {
		message = err.Error()
	}
	writeJSON(w, status, map[string]string{"error": message})
}
