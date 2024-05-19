package service

import (
	"github.com/stretchr/testify/mock"
	"rakia_blog_tt/storage"
)

// MockRepo is a mock implementation of the Repo interface
type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) Create(post storage.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockRepo) GetAll() ([]storage.Post, error) {
	args := m.Called()
	return args.Get(0).([]storage.Post), args.Error(1)
}

func (m *MockRepo) GetByID(id int) (storage.Post, error) {
	args := m.Called(id)
	return args.Get(0).(storage.Post), args.Error(1)
}

func (m *MockRepo) Update(post storage.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *MockRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
