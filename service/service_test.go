package service

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rakia_blog_tt/handler/models"
	"rakia_blog_tt/storage"
)

func loggerMock() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestApplication_CreatePost(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := loggerMock()
	app := New(mockRepo, logger)

	post := models.Post{
		Title:   "Title 1",
		Content: "Content 1",
		Author:  "Author 1",
	}

	dbPost := storage.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}

	mockRepo.On("Create", dbPost).Return(nil)

	err := app.CreatePost(post)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestApplication_GetPosts(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := loggerMock()
	app := New(mockRepo, logger)

	dbPosts := []storage.Post{
		{
			ID:      1,
			Title:   "Title 1",
			Content: "Content 1",
			Author:  "Author 1",
		},
		{
			ID:      2,
			Title:   "Title 2",
			Content: "Content 2",
			Author:  "Author 2",
		},
	}

	mockRepo.On("GetAll").Return(dbPosts, nil)

	posts, err := app.GetPosts()
	require.NoError(t, err)

	expectedPosts := []models.Post{
		{
			ID:      1,
			Title:   "Title 1",
			Content: "Content 1",
			Author:  "Author 1",
		},
		{
			ID:      2,
			Title:   "Title 2",
			Content: "Content 2",
			Author:  "Author 2",
		},
	}

	assert.Equal(t, expectedPosts, posts)
	mockRepo.AssertExpectations(t)
}

func TestApplication_GetPostByID(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := loggerMock()
	app := New(mockRepo, logger)

	dbPost := storage.Post{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
		Author:  "Author 1",
	}

	mockRepo.On("GetByID", 1).Return(dbPost, nil)

	post, err := app.GetPostByID(1)
	require.NoError(t, err)

	expectedPost := models.Post{
		ID:      1,
		Title:   "Title 1",
		Content: "Content 1",
		Author:  "Author 1",
	}

	assert.Equal(t, expectedPost, post)
	mockRepo.AssertExpectations(t)
}

func TestApplication_UpdatePost(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := loggerMock()
	app := New(mockRepo, logger)

	post := models.Post{
		ID:      1,
		Title:   "Updated Title",
		Content: "Updated Content",
		Author:  "Updated Author",
	}

	dbPost := storage.Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
	}

	mockRepo.On("Update", dbPost).Return(nil)

	err := app.UpdatePost(post)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestApplication_DeletePost(t *testing.T) {
	mockRepo := new(MockRepo)
	logger := loggerMock()
	app := New(mockRepo, logger)

	mockRepo.On("Delete", 1).Return(nil)

	err := app.DeletePost(1)
	require.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
