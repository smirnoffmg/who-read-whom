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
	"github.com/what-writers-like/backend/internal/repository"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/service"
	"github.com/what-writers-like/backend/internal/testutils"
)

func setupWriterHandlerRouter(
	t *testing.T,
) (*gin.Engine, repository.WriterRepository, repository.WorkRepository, func()) {
	db, cleanup := testutils.SetupTestDB(t)

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	writerService := service.NewWriterService(writerRepo, workRepo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	writerHandler := handler.NewWriterHandler(writerService)
	router.POST("/writers", writerHandler.Create)
	router.GET("/writers", writerHandler.List)
	router.GET("/writers/:id", writerHandler.GetByID)
	router.PUT("/writers/:id", writerHandler.Update)
	router.DELETE("/writers/:id", writerHandler.Delete)
	return router, writerRepo, workRepo, cleanup
}

func TestWriterHandler_Create(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

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
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

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

	t.Run("invalid birth year", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 0,
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
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, writerRepo, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

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
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/writers/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/writers/999", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWriterHandler_List(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, writerRepo, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charles Dickens", 1812, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))

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
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, writerRepo, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

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

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/writers/invalid", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing name", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/writers/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("writer not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"name":       "Jane Austen",
			"birth_year": 1775,
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/writers/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWriterHandler_Delete(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, writerRepo, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

		req := httptest.NewRequest(http.MethodDelete, "/writers/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodDelete, "/writers/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("writer not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWriterHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodDelete, "/writers/999", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
