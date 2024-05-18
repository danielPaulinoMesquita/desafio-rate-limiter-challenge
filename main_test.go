package main

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	_ "github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Defina uma estrutura para o mock do cliente Redis que incorpora mock.Mock.
type MockRedisClient struct {
	mock.Mock
}

// Implemente o método Get() para o mock do cliente Redis.
func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	args := m.Called(key)
	return args.Get(0).(*redis.StringCmd)
}

func TestHandleRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("API_KEY", "8")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleRequest)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code esperado: %v, mas obteve: %v", http.StatusOK, status)
	}

	expectedBody := "Fazendo Requisição!"
	if rr.Body.String() != expectedBody {
		t.Errorf("Corpo da resposta esperado: %v, mas obteve: %v", expectedBody, rr.Body.String())
	}
}
