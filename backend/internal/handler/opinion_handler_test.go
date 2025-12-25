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

type mockOpinionService struct {
	create      func(uint64, uint64, bool, string, string, *string, *int) (*domain.Opinion, error)
	getByWriter func(uint64) ([]*domain.Opinion, error)
	getByWork   func(uint64) ([]*domain.Opinion, error)
	getOpinion  func(uint64, uint64) (*domain.Opinion, error)
	list        func(int, int) ([]*domain.Opinion, error)
	update      func(uint64, uint64, bool, string, string, *string, *int) error
	delete      func(uint64, uint64) error
}

func (m *mockOpinionService) CreateOpinion(
	writerID, workID uint64,
	sentiment bool,
	quote, source string,
	page *string,
	statementYear *int,
) (*domain.Opinion, error) {
	if m.create != nil {
		return m.create(writerID, workID, sentiment, quote, source, page, statementYear)
	}
	return domain.NewOpinion(writerID, workID, sentiment, quote, source, page, statementYear), nil
}

func (m *mockOpinionService) GetOpinionsByWriter(writerID uint64) ([]*domain.Opinion, error) {
	if m.getByWriter != nil {
		return m.getByWriter(writerID)
	}
	return []*domain.Opinion{}, nil
}

func (m *mockOpinionService) GetOpinionsByWork(workID uint64) ([]*domain.Opinion, error) {
	if m.getByWork != nil {
		return m.getByWork(workID)
	}
	return []*domain.Opinion{}, nil
}

func (m *mockOpinionService) GetOpinion(writerID, workID uint64) (*domain.Opinion, error) {
	if m.getOpinion != nil {
		return m.getOpinion(writerID, workID)
	}
	return domain.NewOpinion(writerID, workID, true, "Quote", "Source", nil, nil), nil
}

func (m *mockOpinionService) ListOpinions(limit, offset int) ([]*domain.Opinion, error) {
	if m.list != nil {
		return m.list(limit, offset)
	}
	return []*domain.Opinion{}, nil
}

func (m *mockOpinionService) UpdateOpinion(
	writerID, workID uint64,
	sentiment bool,
	quote, source string,
	page *string,
	statementYear *int,
) error {
	if m.update != nil {
		return m.update(writerID, workID, sentiment, quote, source, page, statementYear)
	}
	return nil
}

func (m *mockOpinionService) DeleteOpinion(writerID, workID uint64) error {
	if m.delete != nil {
		return m.delete(writerID, workID)
	}
	return nil
}

func setupOpinionHandlerRouter(svc service.OpinionService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	opinionHandler := handler.NewOpinionHandler(svc)
	router.POST("/opinions", opinionHandler.Create)
	router.GET("/opinions", opinionHandler.List)
	router.GET("/opinions/writer/:writer_id", opinionHandler.GetByWriter)
	router.GET("/opinions/work/:work_id", opinionHandler.GetByWork)
	router.GET("/opinions/writer/:writer_id/work/:work_id", opinionHandler.GetByWriterAndWork)
	router.PUT("/opinions/writer/:writer_id/work/:work_id", opinionHandler.Update)
	router.DELETE("/opinions/writer/:writer_id/work/:work_id", opinionHandler.Delete)
	return router
}

func TestOpinionHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{}
		router := setupOpinionHandlerRouter(svc)

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
		svc := &mockOpinionService{}
		router := setupOpinionHandlerRouter(svc)

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

	t.Run("service error", func(t *testing.T) {
		svc := &mockOpinionService{
			create: func(uint64, uint64, bool, string, string, *string, *int) (*domain.Opinion, error) {
				return nil, assert.AnError
			},
		}
		router := setupOpinionHandlerRouter(svc)

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

func TestOpinionHandler_GetByWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			getByWriter: func(uint64) ([]*domain.Opinion, error) {
				return []*domain.Opinion{
					domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
				}, nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

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
		svc := &mockOpinionService{}
		router := setupOpinionHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/invalid", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestOpinionHandler_GetByWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			getByWork: func(uint64) ([]*domain.Opinion, error) {
				return []*domain.Opinion{
					domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
				}, nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

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
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			getOpinion: func(writerID, workID uint64) (*domain.Opinion, error) {
				return domain.NewOpinion(writerID, workID, true, "Quote", "Source", nil, nil), nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

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
		svc := &mockOpinionService{
			getOpinion: func(uint64, uint64) (*domain.Opinion, error) {
				return nil, assert.AnError
			},
		}
		router := setupOpinionHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodGet, "/opinions/writer/2/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestOpinionHandler_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			list: func(int, int) ([]*domain.Opinion, error) {
				return []*domain.Opinion{
					domain.NewOpinion(2, 1, true, "Quote 1", "Source 1", nil, nil),
					domain.NewOpinion(3, 2, false, "Quote 2", "Source 2", nil, nil),
				}, nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

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
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			update: func(uint64, uint64, bool, string, string, *string, *int) error {
				return nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

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
		svc := &mockOpinionService{}
		router := setupOpinionHandlerRouter(svc)

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
	t.Run("success", func(t *testing.T) {
		svc := &mockOpinionService{
			delete: func(uint64, uint64) error {
				return nil
			},
		}
		router := setupOpinionHandlerRouter(svc)

		req := httptest.NewRequest(http.MethodDelete, "/opinions/writer/2/work/1", http.NoBody)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}
