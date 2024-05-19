package handler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"rakia_blog_tt/handler/middleware"
)

func NewRouter(hnd Handler, logger *slog.Logger, metrics middleware.MetricsInterface) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.NewLoggerController(logger, metrics).LoggingMiddleware)
	r.Use(middleware.RemoveTrailingSlash)

	r.Get("/", hnd.DefaultHandler)

	r.Route("/posts", func(r chi.Router) {
		r.Get("/", hnd.GetPosts)          // GET /posts
		r.Post("/", hnd.CreatePost)       // POST /posts
		r.Get("/{id}", hnd.GetPost)       // GET /posts/{id}
		r.Put("/{id}", hnd.UpdatePost)    // PUT /posts/{id}
		r.Delete("/{id}", hnd.DeletePost) // DELETE /posts/{id}
	})

	return r
}
