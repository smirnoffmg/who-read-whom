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

type mockWorkService struct {
	create      func(string, uint64) (*domain.Work, error)
	getByID     func(uint64) (*domain.Work, error)
	getByAuthor func(uint64) ([]*domain.Work, error)
	list        func(int, int) ([]*domain.Work, error)
	update      func(uint64, string, uint64) error
	delete      func(uint64) error
}

func (m *mockWorkService) CreateWork(title string, authorID uint64) (*domain.Work, error) {
	if m.create != nil {
		return m.create(title, authorID)
	}
	return domain.NewWork(1, title, authorID), nil
}

func (m *mockWorkService) GetWork(id uint64) (*domain.Work, error) {
	if m.getByID != nil {
		return m.getByID(id)
	}
	return domain.NewWork(id, "Test Work", 1), nil
}

func (m *mockWorkService) GetWorksByAuthor(authorID uint64) ([]*domain.Work, error) {
	if m.getByAuthor != nil {
		return m.getByAuthor(authorID)
	}
	return []*domain.Work{}, nil
}

func (m *mockWorkService) ListWorks(limit, offset int) ([]*domain.Work, error) {
	if m.list != nil {
		return m.list(limit, offset)
	}
	return []*domain.Work{}, nil
}

func (m *mockWorkService) UpdateWork(id uint64, title string, authorID uint64) error {
	if m.update != nil {
		return m.update(id, title, authorID)
	}
	return nil
}

func (m *mockWorkService) DeleteWork(id uint64) error {
	if m.delete != nil {
		return m.delete(id)
	}
	return nil
}

func setupWorkHandlerRouter(svc service.WorkService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	workHandler := handler.NewWorkHandler(svc)
	router.POST("/works", workHandler.Create)
	router.GET("/works", workHandler.List)
	router.GET("/works/:id", workHandler.GetByID)
	router.GET("/works/author/:author_id", workHandler.GetByAuthor)
	router.PUT("/works/:id", workHandler.Update)
	router.DELETE("/works/:id", workHandler.Delete)
	return router
}

func TestWorkHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{}
		router := setupWorkHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"title":     "Pride and Prejudice",
			"author_id": 1,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/works", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Pride and Prejudice", response["title"])
	})

	t.Run("missing title", func(t *testing.T) {
		svc := &mockWorkService{}
		router := setupWorkHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"author_id": 1,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/works", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockWorkService{
			create: func(string, uint64) (*domain.Work, error) {
				return nil, assert.AnError
			},
		}
		router := setupWorkHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"title":     "Pride and Prejudice",
			"author_id": 1,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/works", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWorkHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{
			getByID: func(id uint64) (*domain.Work, error) {
				return domain.NewWork(id, "Pride and Prejudice", 1), nil
			},
		}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, "Pride and Prejudice", response["title"])
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockWorkService{}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockWorkService{
			getByID: func(uint64) (*domain.Work, error) {
				return nil, assert.AnError
			},
		}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works/999", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWorkHandler_GetByAuthor(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{
			getByAuthor: func(uint64) ([]*domain.Work, error) {
				return []*domain.Work{
					domain.NewWork(1, "Pride and Prejudice", 1),
					domain.NewWork(2, "Sense and Sensibility", 1),
				}, nil
			},
		}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works/author/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
	})

	t.Run("invalid author_id", func(t *testing.T) {
		svc := &mockWorkService{}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works/author/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWorkHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{
			list: func(int, int) ([]*domain.Work, error) {
				return []*domain.Work{
					domain.NewWork(1, "Pride and Prejudice", 1),
					domain.NewWork(2, "Sense and Sensibility", 1),
				}, nil
			},
		}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/works", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
	})
}

func TestWorkHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{
			update: func(uint64, string, uint64) error {
				return nil
			},
		}
		router := setupWorkHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"title":     "Pride and Prejudice (Revised)",
			"author_id": 1,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/works/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing title", func(t *testing.T) {
		svc := &mockWorkService{}
		router := setupWorkHandlerRouter(svc)

		reqBody := map[string]interface{}{
			"author_id": 1,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/works/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWorkHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockWorkService{
			delete: func(uint64) error {
				return nil
			},
		}
		router := setupWorkHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/works/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
