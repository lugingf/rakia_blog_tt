package service

import (
	"log/slog"

	"github.com/pkg/errors"

	"rakia_blog_tt/handler/models"
	"rakia_blog_tt/storage"
)

var ErrPostNotFound = errors.New("post not found")

func New(repo Repo, logger *slog.Logger) *Application {
	return &Application{
		repository: repo,
		logger:     logger,
	}
}

type Application struct {
	repository Repo
	logger     *slog.Logger
}

type Repo interface {
	Create(post storage.Post) error
	GetAll() ([]storage.Post, error)
	GetByID(id int) (storage.Post, error)
	Update(post storage.Post) error
	Delete(id int) error
}

// CreatePost adds a new post
func (app *Application) CreatePost(post models.Post) error {
	dbPost := storage.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}

	app.logger.Debug("Creating a new post")

	return app.repository.Create(dbPost)
}

// GetPosts retrieves all posts
func (app *Application) GetPosts() ([]models.Post, error) {
	app.logger.Debug("Retrieving all posts")

	dbPosts, err := app.repository.GetAll()
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	for _, dbPost := range dbPosts {
		posts = append(posts, models.Post{
			ID:      dbPost.ID,
			Title:   dbPost.Title,
			Content: dbPost.Content,
			Author:  dbPost.Author,
		})
	}

	return posts, nil
}

// GetPostByID retrieves a post by its ID
func (app *Application) GetPostByID(id int) (models.Post, error) {
	app.logger.Debug("Retrieving post by ID", slog.Int("id", id))

	dbPost, err := app.repository.GetByID(id)
	if errors.Is(err, storage.ErrPostNotFound) {
		return models.Post{}, ErrPostNotFound
	}

	return models.Post{
		ID:      dbPost.ID,
		Title:   dbPost.Title,
		Content: dbPost.Content,
		Author:  dbPost.Author,
	}, nil
}

// UpdatePost updates an existing post
func (app *Application) UpdatePost(post models.Post) error {
	dbPost := storage.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}

	app.logger.Debug("Updating post", "post_id", post.ID)
	return app.repository.Update(dbPost)
}

// DeletePost deletes a post by its ID
func (app *Application) DeletePost(id int) error {
	app.logger.Debug("Deleting post", slog.Int("id", id))
	return app.repository.Delete(id)
}
