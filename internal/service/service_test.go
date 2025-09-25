package service

import (
	"context"
	"encoding/csv"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) CreateItem(ctx context.Context, item dto.CreateItem) (*model.Item, error) {
	args := m.Called(ctx, item)
	return args.Get(0).(*model.Item), args.Error(1)
}

func (m *mockStorage) GetAllItems(ctx context.Context, params dto.GetItemsParams) ([]model.Item, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]model.Item), args.Error(1)
}

func (m *mockStorage) UpdateItem(ctx context.Context, id int, item dto.UpdateItem) error {
	args := m.Called(ctx, id, item)
	return args.Error(0)
}

func (m *mockStorage) DeleteItem(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockStorage) GetAggregated(ctx context.Context, from, to string) ([]model.Item, error) {
	args := m.Called(ctx, from, to)
	return args.Get(0).([]model.Item), args.Error(1)
}

func TestNew(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	assert.NotNil(t, s)
	assert.Equal(t, storage, s.storage)
	assert.NotEmpty(t, s.folderName)
}

func TestCreateItem(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	item := dto.CreateItem{Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test"}
	expected := &model.Item{ID: 1, Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test", CreatedAt: time.Now()}
	storage.On("CreateItem", ctx, item).Return(expected, nil)
	result, err := s.CreateItem(ctx, item)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	storage.AssertExpectations(t)
}

func TestGetAllItems(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	params := dto.GetItemsParams{SortBy: []string{"date"}}
	expected := []model.Item{{ID: 1, Type: "расход", Amount: 50, Date: "2023-01-01", Category: "test", CreatedAt: time.Now()}}
	storage.On("GetAllItems", ctx, params).Return(expected, nil)
	result, err := s.GetAllItems(ctx, params)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	storage.AssertExpectations(t)
}

func TestGetAggregated(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	from, to := "2023-01-01", "2023-12-31"
	expected := []model.Item{{ID: 1, Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test", CreatedAt: time.Now(), Aggregated: model.Aggreagated{Sum: 100, Agvarage: 100.0, Count: 1, Median: 100.0, Percentile_90: 100.0}}}
	storage.On("GetAggregated", ctx, from, to).Return(expected, nil)
	result, err := s.GetAggregated(ctx, from, to)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	storage.AssertExpectations(t)
}

func TestUpdateItem(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	id := 1
	item := dto.UpdateItem{Type: stringPtr("расход")}
	storage.On("UpdateItem", ctx, id, item).Return(nil)
	err := s.UpdateItem(ctx, id, item)
	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestDeleteItem(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	id := 1
	storage.On("DeleteItem", ctx, id).Return(nil)
	err := s.DeleteItem(ctx, id)
	assert.NoError(t, err)
	storage.AssertExpectations(t)
}

func TestCSVAggregated(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	from, to := "2023-01-01", "2023-12-31"
	data := []model.Item{{ID: 1, Type: "доход", Amount: 100, Date: "2023-01-01", Category: "test", CreatedAt: time.Now(), Aggregated: model.Aggreagated{Sum: 100, Agvarage: 100.0, Count: 1, Median: 100.0, Percentile_90: 100.0}}}
	storage.On("GetAggregated", ctx, from, to).Return(data, nil)
	fileName, err := s.CSVAggregated(ctx, from, to)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	defer os.Remove(fileName)
	file, err := os.Open(fileName)
	assert.NoError(t, err)
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, []string{"id", "type", "amount", "date", "category", "created_at", "sum", "avarage", "count", "median", "percentile_90"}, records[0])
	assert.Equal(t, []string{"1", "доход", "100", "2023-01-01", "test", data[0].CreatedAt.String(), "100", "100.00", "1", "100.00", "100.00"}, records[1])
	storage.AssertExpectations(t)
}

func TestCSVAllItems(t *testing.T) {
	storage := &mockStorage{}
	s := New(storage)
	ctx := context.Background()
	params := dto.GetItemsParams{}
	data := []model.Item{{ID: 1, Type: "расход", Amount: 50, Date: "2023-01-01", Category: "test", CreatedAt: time.Now()}}
	storage.On("GetAllItems", ctx, params).Return(data, nil)
	fileName, err := s.CSVAllItems(ctx, params)
	assert.NoError(t, err)
	assert.NotEmpty(t, fileName)
	defer os.Remove(fileName)
	file, err := os.Open(fileName)
	assert.NoError(t, err)
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	assert.NoError(t, err)
	assert.Len(t, records, 2)
	assert.Equal(t, []string{"id", "type", "amount", "date", "category", "created_at"}, records[0])
	assert.Equal(t, []string{"1", "расход", "50", "2023-01-01", "test", data[0].CreatedAt.String()}, records[1])
	storage.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
