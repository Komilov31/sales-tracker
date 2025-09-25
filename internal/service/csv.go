package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/Komilov31/sales-tracker/internal/dto"
	"github.com/Komilov31/sales-tracker/internal/model"
)

func (s *Service) CSVAggregated(ctx context.Context, from, to string) (string, error) {
	file, err := os.CreateTemp(s.folderName, "csv*.csv")
	if err != nil {
		return "", fmt.Errorf("could not create csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	data, err := s.storage.GetAggregated(ctx, from, to)
	if err != nil {
		return "", fmt.Errorf("could not get aggregated data: %w", err)
	}

	records := [][]string{{"id", "type", "amount",
		"date", "category", "created_at", "sum",
		"avarage", "count", "median", "percentile_90",
	}}

	for _, record := range getRecords(data, true) {
		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		return "", fmt.Errorf("could not create csv file: %w", err)
	}

	return file.Name(), nil
}

func (s *Service) CSVAllItems(ctx context.Context, params dto.GetItemsParams) (string, error) {
	file, err := os.CreateTemp(s.folderName, "csv*.csv")
	if err != nil {
		return "", fmt.Errorf("could not create csv file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	data, err := s.storage.GetAllItems(ctx, params)
	if err != nil {
		return "", fmt.Errorf("could not get filtered data: %w", err)
	}

	records := [][]string{{"id", "type", "amount", "date", "category", "created_at"}}

	for _, record := range getRecords(data, false) {
		records = append(records, record)
	}

	if err := writer.WriteAll(records); err != nil {
		return "", fmt.Errorf("could not create csv file: %w", err)
	}

	return file.Name(), nil
}

func getRecords(items []model.Item, isAggregated bool) [][]string {
	records := [][]string{}
	for _, item := range items {
		id := fmt.Sprintf("%d", item.ID)
		amount := fmt.Sprintf("%d", item.Amount)
		createdAt := item.CreatedAt.String()
		record := []string{id, item.Type, amount, item.Date, item.Category, createdAt}

		if isAggregated {
			sum := fmt.Sprintf("%d", item.Aggregated.Sum)
			avarage := fmt.Sprintf("%.2f", item.Aggregated.Agvarage)
			count := fmt.Sprintf("%d", item.Aggregated.Count)
			median := fmt.Sprintf("%.2f", item.Aggregated.Median)
			percentile := fmt.Sprintf("%.2f", item.Aggregated.Percentile_90)
			record = append(record, sum, avarage, count, median, percentile)
		}

		records = append(records, record)
	}

	return records
}
