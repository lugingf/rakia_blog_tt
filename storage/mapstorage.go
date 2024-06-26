package storage

import (
	"encoding/json"
	"log/slog"
	"os"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

var ErrPostNotFound = errors.New("post not found")

type Post struct {
	ID      int64
	Title   string
	Content string
	Author  string
}

// InMemoryPostRepository implements the Repo interface
type InMemoryPostRepository struct {
	data   *sync.Map
	nextID int64 // Primary Key, Autoincrement :)
	logger *slog.Logger
}

// NewInMemoryPostRepository creates a new in-memory post repository
func NewInMemoryPostRepository(logger *slog.Logger) *InMemoryPostRepository {
	db := new(sync.Map)
	return &InMemoryPostRepository{data: db, nextID: 1, logger: logger}
}

func (repo *InMemoryPostRepository) Create(post Post) error {
	post.ID = repo.nextID
	repo.nextID++
	repo.data.Store(strconv.FormatInt(post.ID, 10), post)
	return nil
}

func (repo *InMemoryPostRepository) GetAll() ([]Post, error) {
	posts := []Post{}
	repo.data.Range(func(key, value interface{}) bool {
		posts = append(posts, value.(Post))
		return true
	})
	return posts, nil
}

func (repo *InMemoryPostRepository) GetByID(id int) (Post, error) {
	post, ok := repo.data.Load(strconv.Itoa(id))
	if !ok {
		return Post{}, ErrPostNotFound
	}
	return post.(Post), nil
}

func (repo *InMemoryPostRepository) Update(post Post) error {
	_, ok := repo.data.Load(strconv.FormatInt(post.ID, 10))
	if !ok {
		return ErrPostNotFound
	}

	repo.data.Store(strconv.FormatInt(post.ID, 10), post)

	return nil
}

func (repo *InMemoryPostRepository) Delete(id int) error {
	_, ok := repo.data.Load(strconv.Itoa(id))
	if !ok {
		return ErrPostNotFound
	}

	repo.data.Delete(strconv.Itoa(id))

	return nil
}

// saveToFile saves the current state of the repository to a file
// No usage now. Added just in case. Easy to implement and control via config if we need to persist in-mem DB content
func (repo *InMemoryPostRepository) saveToFile(filename string) error {
	posts, err := repo.GetAll()
	if err != nil {
		return err
	}

	data, err := json.Marshal(posts)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// loadFromFile loads the state of the repository from a file
// No usage now. Added just in case. Easy to implement and control via config if we need to persist in-mem DB content
func (repo *InMemoryPostRepository) loadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var posts []Post
	if err := json.Unmarshal(data, &posts); err != nil {
		return err
	}

	repo.data = &sync.Map{}
	for _, post := range posts {
		repo.data.Store(strconv.FormatInt(post.ID, 10), post)
		if post.ID >= repo.nextID {
			repo.nextID = post.ID + 1
		}
	}

	return nil
}
