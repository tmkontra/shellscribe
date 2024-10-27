package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/olivere/vite"
	"github.com/tmkontra/shellscribe/internal/service"
)

type Config struct {
	Directory string
}

type Server struct {
	config  *Config
	service *service.Service
	vite    *vite.Fragment
}

func NewServer(config *Config) *Server {
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:        os.DirFS("web/dist"),
		IsDev:     true,
		ViteURL:   "http://localhost:5173", // optional: defaults to this
		ViteEntry: "src/main.ts",           // reccomended as highly dependent on your app
	})
	if err != nil {
		panic(err)
	}
	return &Server{
		config:  config,
		service: service.NewService(),
		vite:    viteFragment,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/", s.WebHandler)
	r.Handle("/src/assets/*", http.StripPrefix("/src/assets/", http.FileServerFS(os.DirFS("web/src/assets"))))
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
	fmt.Printf("files: %d\n", len(files))
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
	ctx, cancel := context.WithCancel(context.Background())
	go tail(ctx, c)
	for {
		select {
		case <-r.Context().Done():
			log.Printf("SSE client disconnected")
			cancel()
			return
		case line, ok := <-c:
			if !ok {
				log.Printf("tail channel closed")
				cancel()
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", line)
			flusher.Flush()
		}
	}
}

const indexTemplate = `
<head>
    <meta charset="UTF-8" />
    <title>My Go Application</title>
    {{ .Vite.Tags }}
</head>
<body>
	<div id="app"></div>
</body>
`

func (s *Server) WebHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("name").Parse(indexTemplate))

	pageData := map[string]interface{}{
		"Vite": s.vite,
	}

	if err := t.Execute(w, pageData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
