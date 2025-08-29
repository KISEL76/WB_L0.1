package server

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"wb_test/internal/cache"
	"wb_test/internal/storage"
)

//go:embed ui/* ui/static/*
var embedded embed.FS

type Server struct {
	cache  *cache.Cache
	orders *storage.OrderStorage
	mux    *http.ServeMux
	ui     fs.FS
	static fs.FS
}

func New(c *cache.Cache, store *storage.Storage) *Server {
	subUI, _ := fs.Sub(embedded, "ui")
	subStatic, _ := fs.Sub(embedded, "ui/static")

	s := &Server{
		cache:  c,
		orders: store.Orders(),
		mux:    http.NewServeMux(),
		ui:     subUI,
		static: subStatic,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/orders/", s.handleGetOrder)

	fileSrv := http.FileServer(http.FS(s.static))
	s.mux.Handle("/static/", http.StripPrefix("/static/", fileSrv))

	s.mux.HandleFunc("/", s.handleIndex)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(s.ui, "index.html")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (s *Server) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/orders/")
	if id == "" {
		http.Error(w, "order_id required", http.StatusBadRequest)
		return
	}

	if order, ok := s.cache.Get(id); ok {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(order)
		return
	}

	order, err := s.orders.Get(r.Context(), id)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if order == nil {
		http.Error(w, "error not found", http.StatusNotFound)
		return
	}

	s.cache.Set(order)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(order)
}

func (s *Server) Start(addr string) error {
	log.Printf("HTTP server started on %s", addr)
	return http.ListenAndServe(addr, s.mux)
}
