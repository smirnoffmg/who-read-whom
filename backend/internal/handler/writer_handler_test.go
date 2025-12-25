package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/what-writers-like/backend/internal/domain"
	"github.com/what-writers-like/backend/internal/handler"
	"github.com/what-writers-like/backend/internal/service"
)

type mockWriterService struct {
	create  func(string, int, *int, *string) (*domain.Writer, error)
	getByID func(uint64) (*domain.Writer, error)
	list    func(int, int) ([]*domain.Writer, error)
	update  func(uint64, string, int, *int, *string) error
	delete  func(uint64) error
}

func (m *mockWriterService) CreateWriter(
	name string,
	birthYear int,
	deathYear *int,
	bio *string,
) (*domain.Writer, error) {
	if m.create != nil {
		return m.create(name, birthYear, deathYear, bio)
	}
	return domain.NewWriter(1, name, birthYear, deathYear, bio), nil
}

func (m *mockWriterService) GetWriter(id uint64) (*domain.Writer, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return domain.NewWriter(id, "Test Writer", 1800, nil, nil), nil
}

func (m *mockWriterService) ListWriters(limit, offset int) ([]*domain.Writer, error) {
	if m.list != nil {
		return m.list(limit, offset)
	}
	return []*domain.Writer{}, nil
}

func (m *mockWriterService) UpdateWriter(id uint64, name string, birthYear int, deathYear *int, bio *string) error {
	if m.update != nil {
		return m.update(id, name, birthYear, deathYear, bio)
	}
	return nil
}

func (m *mockWriterService) DeleteWriter(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}

func setupWriterHandlerRouter(service service.WriterService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	writerHandler := handler.NewWriterHandler(service)
	router.POST("/writers", writerHandler.Create)
	router.GET("/writers", writerHandler.List)
	router.GET("/writers/:id", writerHandler.GetByID)
	router.PUT("/writers/:id", writerHandler.Update)
	router.DELETE("/writers/:id", writerHandler.Delete)
	return router
}

func TestWriterHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWriterService{}
		router := setupWriterHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/writers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Jane Austen", response["name"])
	})

	t.Run("missing name", func(t *testing.T) {
		svc := &mockWriterService{}
		router := setupWriterHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/writers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockWriterService{
			create: func(string, int, *int, *string) (*domain.Writer, error) {
				return nil, assert.AnError
			},
		}
		router := setupWriterHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/writers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWriterHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWriterService{
			getByID: func(id uint64) (*domain.Writer, error) {
				return domain.NewWriter(id, "Jane Austen", 1775, nil, nil), nil
			},
		}
		router := setupWriterHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/writers/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Jane Austen", response["name"])
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockWriterService{}
		router := setupWriterHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/writers/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockWriterService{
			getByID: func(uint64) (*domain.Writer, error) {
				return nil, assert.AnError
			},
		}
		router := setupWriterHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/writers/999", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWriterHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWriterService{
			list: func(int, int) ([]*domain.Writer, error) {
				return []*domain.Writer{
					domain.NewWriter(1, "Jane Austen", 1775, nil, nil),
					domain.NewWriter(2, "Charles Dickens", 1812, nil, nil),
				}, nil
			},
		}
		router := setupWriterHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/writers", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
	})
}

func TestWriterHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWriterService{
			update: func(uint64, string, int, *int, *string) error {
				return nil
			},
		}
		router := setupWriterHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/writers/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestWriterHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWriterService{
			delete: func(uint64) error {
				return nil
			},
		}
		router := setupWriterHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/writers/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
