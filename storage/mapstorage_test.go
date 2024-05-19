package storage

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loggerMock() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
}

func TestInMemoryPostRepository(t *testing.T) {
	logger := loggerMock()
	repo := NewInMemoryPostRepository(logger)

	t.Run("Create Post", func(t *testing.T) {
		post := Post{Title: "Title 1", Content: "Content 1", Author: "Author 1"}
		err := repo.Create(post)
		require.NoError(t, err)

		// Since we don't have the post ID directly after creation, we'll retrieve all posts
		posts, err := repo.GetAll()
		require.NoError(t, err)
		require.Len(t, posts, 1)

		retrievedPost := posts[0]
		assert.Equal(t, post.Title, retrievedPost.Title)
		assert.Equal(t, post.Content, retrievedPost.Content)
		assert.Equal(t, post.Author, retrievedPost.Author)
		assert.NotZero(t, int(retrievedPost.ID))
	})

	t.Run("Get All Posts", func(t *testing.T) {
		posts, err := repo.GetAll()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)
	})

	t.Run("Get Post By ID", func(t *testing.T) {
		post := Post{Title: "Title 2", Content: "Content 2", Author: "Author 2"}
		err := repo.Create(post)
		require.NoError(t, err)

		// Since we don't have the post ID directly after creation, we'll retrieve all posts
		posts, err := repo.GetAll()
		require.NoError(t, err)
		require.Len(t, posts, 2) // Should be 2 posts now

		retrievedPost := posts[1] // The second post
		retrievedPostByID, err := repo.GetByID(int(retrievedPost.ID))
		require.NoError(t, err)
		assert.Equal(t, retrievedPost, retrievedPostByID)
	})

	t.Run("Update Post", func(t *testing.T) {
		post := Post{Title: "Title 3", Content: "Content 3", Author: "Author 3"}
		err := repo.Create(post)
		require.NoError(t, err)

		// Retrieve all posts to get the ID of the last inserted post
		posts, err := repo.GetAll()
		require.NoError(t, err)
		require.Len(t, posts, 3)

		// Update the last post
		retrievedPost := posts[2]
		retrievedPost.Title = "Updated Title 3"
		err = repo.Update(retrievedPost)
		require.NoError(t, err)

		// Verify update
		updatedPost, err := repo.GetByID(int(retrievedPost.ID))
		require.NoError(t, err)
		assert.Equal(t, "Updated Title 3", updatedPost.Title)
	})

	t.Run("Delete Post", func(t *testing.T) {
		post := Post{Title: "Title 4", Content: "Content 4", Author: "Author 4"}
		err := repo.Create(post)
		require.NoError(t, err)

		// Retrieve all posts to get the ID of the last inserted post
		posts, err := repo.GetAll()
		require.NoError(t, err)
		require.Len(t, posts, 4)

		// Delete the last post
		retrievedPost := posts[3]
		err = repo.Delete(int(retrievedPost.ID))
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetByID(int(retrievedPost.ID))
		assert.Error(t, err)
	})

	t.Run("Save To File and Load From File", func(t *testing.T) {
		filename := "test_posts.json"
		defer os.Remove(filename) // Clean up after test

		err := repo.saveToFile(filename)
		require.NoError(t, err)

		newRepo := NewInMemoryPostRepository(logger)
		err = newRepo.loadFromFile(filename)
		require.NoError(t, err)

		posts, err := newRepo.GetAll()
		require.NoError(t, err)
		assert.NotEmpty(t, posts)
	})
}
