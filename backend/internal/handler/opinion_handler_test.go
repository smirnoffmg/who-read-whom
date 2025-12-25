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

func setupOpinionHandlerRouter(
	t *testing.T,
) (*gin.Engine, repository.OpinionRepository, repository.WriterRepository, repository.WorkRepository, func()) {
	db, cleanup := testutils.SetupTestDB(t)

	opinionRepo := gorm.NewOpinionRepository(db)
	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionService := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	opinionHandler := handler.NewOpinionHandler(opinionService)
	router.POST("/opinions", opinionHandler.Create)
	router.GET("/opinions", opinionHandler.List)
	router.GET("/opinions/writer/:writer_id", opinionHandler.GetByWriter)
	router.GET("/opinions/work/:work_id", opinionHandler.GetByWork)
	router.GET("/opinions/writer/:writer_id/work/:work_id", opinionHandler.GetByWriterAndWork)
	router.PUT("/opinions/writer/:writer_id/work/:work_id", opinionHandler.Update)
	router.DELETE("/opinions/writer/:writer_id/work/:work_id", opinionHandler.Delete)
	return router, opinionRepo, writerRepo, workRepo, cleanup
}

func TestOpinionHandler_Create(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, _, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		// Create test data
		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		reqBody := map[string]interface{}{
			"writer_id": 2,
			"work_id":   1,
			"sentiment": true,
			"quote":     "A delightful novel",
			"source":    "Personal Letters",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/opinions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, true, response["sentiment"])
		assert.Equal(t, "A delightful novel", response["quote"])
	})

	t.Run("missing quote", func(t *testing.T) {
		t.Parallel()
		router, _, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		// Create test data
		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		reqBody := map[string]interface{}{
			"writer_id": 2,
			"work_id":   1,
			"sentiment": true,
			"source":    "Personal Letters",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/opinions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("writer not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, _, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		reqBody := map[string]interface{}{
			"writer_id": 2,
			"work_id":   1,
			"sentiment": true,
			"quote":     "A delightful novel",
			"source":    "Personal Letters",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/opinions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func setupTestOpinionData(
	t *testing.T,
	writerRepo repository.WriterRepository,
	workRepo repository.WorkRepository,
	opinionRepo repository.OpinionRepository,
) {
	writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
	writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
	require.NoError(t, writerRepo.Create(writer1))
	require.NoError(t, writerRepo.Create(writer2))
	work := domain.NewWork(1, "Pride and Prejudice", 1)
	require.NoError(t, workRepo.Create(work))
	opinion := domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil)
	require.NoError(t, opinionRepo.Create(opinion))
}

func TestOpinionHandler_GetByWriter(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		setupTestOpinionData(t, writerRepo, workRepo, opinionRepo)

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/2", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})

	t.Run("invalid writer_id", func(t *testing.T) {
		t.Parallel()
		router, _, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		// Create test data
		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOpinionHandler_GetByWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		setupTestOpinionData(t, writerRepo, workRepo, opinionRepo)

		req := httptest.NewRequest(http.MethodGet, "/opinions/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(response), 1)
	})
}

func TestOpinionHandler_GetByWriterAndWork(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))
		opinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
		require.NoError(t, opinionRepo.Create(opinion))

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/2/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Equal(t, uint64(2), uint64(response["writer_id"].(float64)))
		assert.Equal(t, uint64(1), uint64(response["work_id"].(float64)))
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		router, _, _, _, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/2/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestOpinionHandler_List(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		writer3 := domain.NewWriter(3, "Charles Dickens", 1812, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		require.NoError(t, writerRepo.Create(writer3))
		work1 := domain.NewWork(1, "Pride and Prejudice", 1)
		work2 := domain.NewWork(2, "Jane Eyre", 2)
		require.NoError(t, workRepo.Create(work1))
		require.NoError(t, workRepo.Create(work2))
		opinion1 := domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil)
		opinion2 := domain.NewOpinion(3, 2, false, "Quote 2", "Source 2", nil, nil)
		require.NoError(t, opinionRepo.Create(opinion1))
		require.NoError(t, opinionRepo.Create(opinion2))

		req := httptest.NewRequest(http.MethodGet, "/opinions", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		assert.Len(t, response, 2)
	})
}

func TestOpinionHandler_Update(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))
		opinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
		require.NoError(t, opinionRepo.Create(opinion))

		reqBody := map[string]interface{}{
			"sentiment": true,
			"quote":     "Updated quote",
			"source":    "Updated source",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/opinions/writer/2/work/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing quote", func(t *testing.T) {
		t.Parallel()
		router, _, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		// Create test data
		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))

		reqBody := map[string]interface{}{
			"sentiment": true,
			"source":    "Source",
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/opinions/writer/2/work/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOpinionHandler_Delete(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		router, opinionRepo, writerRepo, workRepo, cleanup := setupOpinionHandlerRouter(t)
		defer cleanup()

		writer1 := domain.NewWriter(1, "Jane Austen", 1775, nil, nil)
		writer2 := domain.NewWriter(2, "Charlotte Bronte", 1816, nil, nil)
		require.NoError(t, writerRepo.Create(writer1))
		require.NoError(t, writerRepo.Create(writer2))
		work := domain.NewWork(1, "Pride and Prejudice", 1)
		require.NoError(t, workRepo.Create(work))
		opinion := domain.NewOpinion(2, 1, true, "Quote", "Source", nil, nil)
		require.NoError(t, opinionRepo.Create(opinion))

		req := httptest.NewRequest(http.MethodDelete, "/opinions/writer/2/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
