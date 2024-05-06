package mocks

import "github.com/stretchr/testify/mock"

type MockTransactionStorage struct {
	mock.Mock
}

func (m *MockTransactionStorage) Save(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockTransactionStorage) Find(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockTransactionStorage) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockTransactionStorage) Update(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockTransactionStorage) GetAll() interface{} {
	args := m.Called()
	return args.Get(0)
}
