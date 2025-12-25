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

func setupWorkHandlerRouter(
	t *testing.T,
) (*gin.Engine, repository.WorkRepository, repository.WriterRepository, func()) {
	db, cleanup := testutils.SetupTestDB(t)

	workRepo := gorm.NewWorkRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workService := service.NewWorkService(workRepo, writerRepo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	workHandler := handler.NewWorkHandler(workService)
	router.POST("/works", workHandler.Create)
	router.GET("/works", workHandler.List)
	router.GET("/works/:id", workHandler.GetByID)
	router.GET("/works/author/:author_id", workHandler.GetByAuthor)
	router.PUT("/works/:id", workHandler.Update)
	router.DELETE("/works/:id", workHandler.Delete)
	return router, workRepo, writerRepo, cleanup
}

func TestWorkHandler_Create(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, _, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		// Create writer first
		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))

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
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

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

	t.Run("author not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"title":     "Pride and Prejudice",
			"author_id": 999,
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
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, workRepo, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		// Create writer and work first
		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

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
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/works/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/works/999", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestWorkHandler_GetByAuthor(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, workRepo, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work1 := domain.NewWork(1, "Pride and Prejudice", 1)
		work2 := domain.NewWork(2, "Sense and Sensibility", 1)
		require.NoError(t, workRepo.Create(work1))
		require.NoError(t, workRepo.Create(work2))

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
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/works/author/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestWorkHandler_List(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, workRepo, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work1 := domain.NewWork(1, "Pride and Prejudice", 1)
		work2 := domain.NewWork(2, "Sense and Sensibility", 1)
		require.NoError(t, workRepo.Create(work1))
		require.NoError(t, workRepo.Create(work2))

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
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, workRepo, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

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
		t.Parallel()
		router, _, _, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

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
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, workRepo, writerRepo, cleanup := setupWorkHandlerRouter(t)
		defer cleanup()

		writer := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		require.NoError(t, writerRepo.Create(writer))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		req := httptest.NewRequest(http.MethodDelete, "/works/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
