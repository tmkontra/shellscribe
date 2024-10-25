package server

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/tmkontra/shellscribe/internal/service"
)

type Config struct {
	Directory string
}

type Server struct {
	config  *Config
	service *service.Service
}

func NewServer(config *Config) *Server {
	return &Server{
		config:  config,
		service: service.NewService(),
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/index", s.IndexHandler)
	r.Get("/tail/{id}", s.TailHandler)
	r.ServeHTTP(w, req)
}

type LogFile struct {
	Id        string
	Name      string
	CreatedAt time.Time
}

type IndexResponse struct {
	Logs []LogFile
}

func (s *Server) IndexHandler(w http.ResponseWriter, r *http.Request) {
	files, err := s.service.ListFiles(s.config.Directory)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		return
	}
	resp := IndexResponse{}
	for _, file := range files {
		resp.Logs = append(resp.Logs, LogFile{
			Id:        url.PathEscape(file.Id),
			Name:      file.Cmd,
			CreatedAt: file.StartTimestamp,
		})
	}
	render.Status(r, http.StatusOK)
	render.JSON(w, r, resp)
}

func (s *Server) TailHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := url.PathUnescape(idParam)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, nil)
		return
	}
	tail, err := s.service.TailFile(id + "/output")
	if err != nil {
		fmt.Printf("tail error: %s\n", err)
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, nil)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		render.Status(r, http.StatusInternalServerError)
		render.PlainText(w, r, "Streaming unsupported.")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	c := make(chan string)
	go tail(c)
	for {
		select {
		case <-r.Context().Done():
			log.Printf("SSE client disconnected")
			return
		case line := <-c:
			w.Write([]byte(line))
			flusher.Flush()
		}
	}
}
