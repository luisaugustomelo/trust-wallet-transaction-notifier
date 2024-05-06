package mocks

import "github.com/stretchr/testify/mock"

type MockSubscriptionStorage struct {
	mock.Mock
}

func (m *MockSubscriptionStorage) Save(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockSubscriptionStorage) Find(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

func (m *MockSubscriptionStorage) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockSubscriptionStorage) Update(key string, value interface{}) {
	m.Called(key, value)
}

func (m *MockSubscriptionStorage) GetAll() interface{} {
	args := m.Called()
	return args.Get(0)
}
