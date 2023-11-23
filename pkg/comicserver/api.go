package comicserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/klaital/comics/pkg/datalayer"
	"github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	storage *datalayer.ComicDataSource
	l       *logrus.Entry
	mux     *chi.Mux
}

func New(storage *datalayer.ComicDataSource, logger *logrus.Entry) *Server {
	srv := Server{
		storage: storage,
		l:       logger,
	}
	// Ensure a default logger is attached
	if logger == nil {
		srv.l = logrus.NewEntry(logrus.New())
	}

	srv.mux = chi.NewRouter()
	srv.mux.Use(middleware.Timeout(10 * time.Second))
	srv.mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost", "https://comics.klaital.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"link"},
		AllowCredentials: false,
		MaxAge:           300,
		Debug:            false,
	}))

	// Configure the exposed routes
	srv.mux.Get("/api/comics", srv.GetAllComics)
	srv.mux.Post("/api/comics/{comicID}/rss", srv.RefreshRssFeed)
	srv.mux.Put("/api/comics/{comicID}/read", srv.MarkComicRead)
	srv.mux.Put("/api/comics/{comicID}/rss/{rssItemId}/read", srv.MarkItemRead)
	return &srv
}

func (srv *Server) ServeHTTP() {

}
