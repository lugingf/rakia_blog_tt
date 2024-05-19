package handler

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/strfmt"
	"log/slog"
	"net/http"
	"strconv"

	"rakia_blog_tt/handler/models"
	"rakia_blog_tt/service"
)

func New(app *service.Application, logger *slog.Logger) Handler {
	return Handler{
		service: app,
		logger:  logger,
	}
}

type Handler struct {
	service *service.Application
	logger  *slog.Logger
}

func (h *Handler) DefaultHandler(w http.ResponseWriter, req *http.Request) {
	respText := map[string]string{"response": "test BE response"}
	rj, _ := json.Marshal(respText)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rj) // nolint:errcheck
}

func (h *Handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.GetPosts()
	if err != nil {
		h.logger.Error("failed to get posts", "error", err)
		http.Error(w, "Failed to retrieve posts", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts) // nolint:errcheck
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		h.logger.Error("post decode failed", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if err := post.Validate(strfmt.NewFormats()); err != nil {
		h.logger.Error("invalid post format", "error", err)
		http.Error(w, "Invalid post format", http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePost(post); err != nil {
		h.logger.Error("failed to create post", "error", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	post, err := h.service.GetPostByID(id)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			h.logger.Error("failed to get post", "error", err)
			http.Error(w, "Failed to retrieve post", http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(post) // nolint:errcheck
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	var post models.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		h.logger.Error("post decode failed", "error", err)
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	post.ID = int64(id)
	if err := h.service.UpdatePost(post); err != nil {
		h.logger.Error("failed to update post", "error", err)
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeletePost(id); err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			h.logger.Error("failed to delete post", "error", err)
			http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
