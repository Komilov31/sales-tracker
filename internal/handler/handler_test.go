package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockTrackerService struct {
	mock.Mock
}

func (m *mockTrackerService) CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *mockTrackerService) GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *mockTrackerService) GetAggregated(ctx context.Context, from, to string) ([]model.Item, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *mockTrackerService) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	args := m.Called(ctx, id, item)
	return args.Error(0)
}

func (m *mockTrackerService) DeleteItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTrackerService) CSVAggregated(ctx context.Context, from, to string) (string, error) {
	args := m.Called(ctx, from, to)
	return args.String(0), args.Error(1)
}

func (m *mockTrackerService) CSVAllItems(ctx context.Context, params dto.GetItemsParams) (string, error) {
	args := m.Called(ctx, params)
	return args.String(0), args.Error(1)
}

func TestCreateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		item := dto.CreateItem{Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test"}
		expected := &model.Item{ID: 1, Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test"}
		mockService.On("CreateItem", mock.Anything, item).Return(expected, nil)

		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.CreateItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response dto.ItemWithoutAggregated
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, expected.ID, response.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.CreateItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateItem")
	})

	t.Run("validation error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		item := dto.CreateItem{Type: "", Amount: 100, Date: "2023-01-01", Category: "test"}
		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.CreateItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CreateItem")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		item := dto.CreateItem{Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test"}
		mockService.On("CreateItem", mock.Anything, item).Return((*model.Item)(nil), assert.AnError)

		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPost, "/items", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.CreateItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAllItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		params := dto.GetItemsParams{SortBy: []string{"date"}}
		expected := []model.Item{{ID: 1, Type: "расход", Amount: 50, Date: "2023-01-01", Category: "test"}}
		mockService.On("GetAllItems", mock.Anything, params).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/items?sort_by=date", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAllItems(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []dto.ItemWithoutAggregated
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 1)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid sort_by", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodGet, "/items?sort_by=invalid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAllItems(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetAllItems")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		params := dto.GetItemsParams{}
		mockService.On("GetAllItems", mock.Anything, params).Return([]model.Item(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/items", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAllItems(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAggregated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		from, to := "2023-01-01", "2023-12-31"
		expected := []model.Item{{ID: 1, Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test"}}
		mockService.On("GetAggregated", mock.Anything, from, to).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/analytics?from=2023-01-01&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregated(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.Item
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Len(t, response, 1)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid date", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodGet, "/analytics?from=invalid&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregated(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetAggregated")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		from, to := "2023-01-01", "2023-12-31"
		mockService.On("GetAggregated", mock.Anything, from, to).Return([]model.Item(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/analytics?from=2023-01-01&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregated(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUpdateItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		id := 1
		item := dto.UpdateItem{Type: stringPtr("расход")}
		mockService.On("UpdateItem", mock.Anything, id, item).Return(nil)

		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPut, "/items/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodPut, "/items/invalid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.UpdateItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("invalid json", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodPut, "/items/1", bytes.NewBufferString("invalid"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		id := 1
		item := dto.UpdateItem{}
		mockService.On("UpdateItem", mock.Anything, id, item).Return(assert.AnError)

		body, _ := json.Marshal(item)
		req := httptest.NewRequest(http.MethodPut, "/items/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.UpdateItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestDeleteItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		id := 1
		mockService.On("DeleteItem", mock.Anything, id).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteItem(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodDelete, "/items/invalid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "invalid"}}

		handler.DeleteItem(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "DeleteItem")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		id := 1
		mockService.On("DeleteItem", mock.Anything, id).Return(assert.AnError)

		req := httptest.NewRequest(http.MethodDelete, "/items/1", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req
		c.Params = gin.Params{{Key: "id", Value: "1"}}

		handler.DeleteItem(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetAggregatedCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		from, to := "2023-01-01", "2023-12-31"
		tempFile, err := os.CreateTemp("", "test.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tempFile.Name())
		mockService.On("CSVAggregated", mock.Anything, from, to).Return(tempFile.Name(), nil)

		req := httptest.NewRequest(http.MethodGet, "/analytics/csv?from=2023-01-01&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregatedCSV(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid date", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodGet, "/analytics/csv?from=invalid&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregatedCSV(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CSVAggregated")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		from, to := "2023-01-01", "2023-12-31"
		mockService.On("CSVAggregated", mock.Anything, from, to).Return("", assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/analytics/csv?from=2023-01-01&to=2023-12-31", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetAggregatedCSV(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetFilteredCSV(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("success", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		params := dto.GetItemsParams{SortBy: []string{"date"}}
		tempFile, err := os.CreateTemp("", "test.csv")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tempFile.Name())
		mockService.On("CSVAllItems", mock.Anything, params).Return(tempFile.Name(), nil)

		req := httptest.NewRequest(http.MethodGet, "/items/csv?sort_by=date", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetFilteredCSV(c)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid sort_by", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		req := httptest.NewRequest(http.MethodGet, "/items/csv?sort_by=invalid", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetFilteredCSV(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "CSVAllItems")
	})

	t.Run("service error", func(t *testing.T) {
		mockService := &mockTrackerService{}
		handler := New(context.Background(), mockService)
		params := dto.GetItemsParams{}
		mockService.On("CSVAllItems", mock.Anything, params).Return("", assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/items/csv", nil)
		w := httptest.NewRecorder()

		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handler.GetFilteredCSV(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func stringPtr(s string) *string {
	return &s
}
