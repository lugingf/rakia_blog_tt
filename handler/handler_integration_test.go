package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"rakia_blog_tt/handler/models"
	"rakia_blog_tt/service"
	"rakia_blog_tt/storage"
)

type metricsMock struct {
}

func (m *metricsMock) ObserveHTTPDuration(timeSince time.Time, path string, code int, method string) {
}

func loggerMock() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func setupTestServer() *httptest.Server {
	logger := loggerMock()
	postRepo := storage.NewInMemoryPostRepository(logger)
	application := service.New(postRepo, logger)
	hndl := New(application, logger)
	router := NewRouter(hndl, logger, &metricsMock{})

	return httptest.NewServer(router)
}

func TestIntegration_CreatePostHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	post := models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Test Author",
	}
	body, err := json.Marshal(post)
	require.NoError(t, err)

	resp, err := http.Post(server.URL+"/posts", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestIntegration_GetPostsHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a post first
	post := models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Test Author",
	}
	body, err := json.Marshal(post)
	require.NoError(t, err)
	_, err = http.Post(server.URL+"/posts", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)

	resp, err := http.Get(server.URL + "/posts")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var posts []models.Post
	err = json.NewDecoder(resp.Body).Decode(&posts)
	require.NoError(t, err)
	assert.NotEmpty(t, posts)
}

func TestIntegration_GetPostHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a post first
	post := models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Test Author",
	}
	body, err := json.Marshal(post)
	require.NoError(t, err)
	resp, err := http.Post(server.URL+"/posts", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)

	// Parse the created post ID from the response
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// Fetch the created post
	resp, err = http.Get(server.URL + "/posts/1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedPost models.Post
	err = json.NewDecoder(resp.Body).Decode(&retrievedPost)
	require.NoError(t, err)
	assert.Equal(t, post.Title, retrievedPost.Title)
	assert.Equal(t, post.Content, retrievedPost.Content)
	assert.Equal(t, post.Author, retrievedPost.Author)
}

func TestIntegration_UpdatePostHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a post first
	post := models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Test Author",
	}
	body, err := json.Marshal(post)
	require.NoError(t, err)
	resp, err := http.Post(server.URL+"/posts", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	resp.Body.Close()

	// Update the created post
	updatedPost := models.Post{
		Title:   "Updated Test Post",
		Content: "This is an updated test post",
		Author:  "Updated Test Author",
	}
	updatedBody, err := json.Marshal(updatedPost)
	require.NoError(t, err)
	req, err := http.NewRequest(http.MethodPut, server.URL+"/posts/1", bytes.NewBuffer(updatedBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Fetch the updated post
	resp, err = http.Get(server.URL + "/posts/1")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedPost models.Post
	err = json.NewDecoder(resp.Body).Decode(&retrievedPost)
	require.NoError(t, err)
	assert.Equal(t, updatedPost.Title, retrievedPost.Title)
	assert.Equal(t, updatedPost.Content, retrievedPost.Content)
	assert.Equal(t, updatedPost.Author, retrievedPost.Author)
}

func TestIntegration_DeletePostHandler(t *testing.T) {
	server := setupTestServer()
	defer server.Close()

	// Create a post first
	post := models.Post{
		Title:   "Test Post",
		Content: "This is a test post",
		Author:  "Test Author",
	}
	body, err := json.Marshal(post)
	require.NoError(t, err)
	resp, err := http.Post(fmt.Sprintf("%s/posts", server.URL), "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	resp.Body.Close()

	// Delete the created post
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/posts/1", server.URL), nil)
	require.NoError(t, err)
	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Fetch the deleted post
	resp, err = http.Get(fmt.Sprintf("%s/posts/1", server.URL))
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
