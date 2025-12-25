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
	"github.com/what-writers-like/backend/internal/handler"
	"github.com/what-writers-like/backend/internal/repository/gorm"
	"github.com/what-writers-like/backend/internal/service"
	"github.com/what-writers-like/backend/internal/testutils"
)

func setupE2ERouter(t *testing.T) (*gin.Engine, func()) {
	db, cleanup := testutils.SetupTestDB(t)

	writerRepo := gorm.NewWriterRepository(db)
	workRepo := gorm.NewWorkRepository(db)
	opinionRepo := gorm.NewOpinionRepository(db)

	writerService := service.NewWriterService(writerRepo, workRepo)
	workService := service.NewWorkService(workRepo, writerRepo)
	opinionService := service.NewOpinionService(opinionRepo, writerRepo, workRepo)

	writerHandler := handler.NewWriterHandler(writerService)
	workHandler := handler.NewWorkHandler(workService)
	opinionHandler := handler.NewOpinionHandler(opinionService)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router = handler.SetupRouter(writerHandler, workHandler, opinionHandler)

	return router, cleanup
}

func TestE2E_WriterWorkflow(t *testing.T) {
	router, cleanup := setupE2ERouter(t)
	defer cleanup()

	// Create writer
	createReq := map[string]interface{}{
		"name":       "Jane Austen",
		"birth_year": 1775,
	}
	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	writerID := uint64(createResp["id"].(float64))

	// Get writer
	req = httptest.NewRequest(http.MethodGet, "/api/v1/writers/1", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var getResp map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	require.NoError(t, err)
	assert.Equal(t, "Jane Austen", getResp["name"])

	// List writers
	req = httptest.NewRequest(http.MethodGet, "/api/v1/writers", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var listResp []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp), 1)

	// Update writer
	bio := "English novelist"
	updateReq := map[string]interface{}{
		"name":       "Jane Austen",
		"birth_year": 1775,
		"bio":        bio,
	}
	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/writers/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Verify update
	req = httptest.NewRequest(http.MethodGet, "/api/v1/writers/1", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	err = json.Unmarshal(w.Body.Bytes(), &getResp)
	require.NoError(t, err)
	assert.Equal(t, bio, getResp["bio"])

	_ = writerID
}

func TestE2E_WorkWorkflow(t *testing.T) {
	router, cleanup := setupE2ERouter(t)
	defer cleanup()

	// Create writer first
	createWriterReq := map[string]interface{}{
		"name":       "Jane Austen",
		"birth_year": 1775,
	}
	body, _ := json.Marshal(createWriterReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Create work
	createWorkReq := map[string]interface{}{
		"title":     "Pride and Prejudice",
		"author_id": 1,
	}
	body, _ = json.Marshal(createWorkReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/works", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	assert.Equal(t, "Pride and Prejudice", createResp["title"])

	// Get work by author
	req = httptest.NewRequest(http.MethodGet, "/api/v1/works/author/1", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var listResp []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp), 1)
}

func TestE2E_OpinionWorkflow(t *testing.T) {
	router, cleanup := setupE2ERouter(t)
	defer cleanup()

	// Create two writers
	createWriter1Req := map[string]interface{}{
		"name":       "Jane Austen",
		"birth_year": 1775,
	}
	body, _ := json.Marshal(createWriter1Req)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	createWriter2Req := map[string]interface{}{
		"name":       "Charlotte Bronte",
		"birth_year": 1816,
	}
	body, _ = json.Marshal(createWriter2Req)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/writers", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Create work
	createWorkReq := map[string]interface{}{
		"title":     "Pride and Prejudice",
		"author_id": 1,
	}
	body, _ = json.Marshal(createWorkReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/works", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	// Create opinion
	createOpinionReq := map[string]interface{}{
		"writer_id": 2,
		"work_id":   1,
		"sentiment": true,
		"quote":     "A delightful novel",
		"source":    "Personal Letters",
	}
	body, _ = json.Marshal(createOpinionReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/opinions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)
	var createResp map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &createResp)
	require.NoError(t, err)
	assert.Equal(t, true, createResp["sentiment"])
	assert.Equal(t, "A delightful novel", createResp["quote"])

	// Get opinions by writer
	req = httptest.NewRequest(http.MethodGet, "/api/v1/opinions/writer/2", http.NoBody)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	var listResp []map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &listResp)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(listResp), 1)

	// Try to create opinion about own work (should fail)
	createOwnOpinionReq := map[string]interface{}{
		"writer_id": 1,
		"work_id":   1,
		"sentiment": true,
		"quote":     "My own work",
		"source":    "Personal",
	}
	body, _ = json.Marshal(createOwnOpinionReq)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/opinions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
