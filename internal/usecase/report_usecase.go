package usecase

import (
	"context"
	"kasir-api/internal/models"
	"kasir-api/internal/repository"
	"net/http"
	"time"
)

type ReportUseCase struct {
	TransactionRepository *repository.TransactionRepository
}

func NewReportUseCase(transactionRepository *repository.TransactionRepository) *ReportUseCase {
	return &ReportUseCase{
		TransactionRepository: transactionRepository,
	}
}

func (uc *ReportUseCase) Today(ctx context.Context) (*models.ReportResponse, error) {
	now := time.Now()
	loc := now.Location()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour)

	return uc.getReport(ctx, start, end)
}

func (uc *ReportUseCase) ByDateRange(ctx context.Context, startDate, endDate string) (*models.ReportResponse, error) {
	if startDate == "" || endDate == "" {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: "start_date and end_date are required",
		}
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: "invalid start_date format, use YYYY-MM-DD",
		}
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: "invalid end_date format, use YYYY-MM-DD",
		}
	}

	if end.Before(start) {
		return nil, &models.Error{
			Status:  http.StatusBadRequest,
			Message: "end_date must be after or equal to start_date",
		}
	}

	end = end.Add(24 * time.Hour)

	return uc.getReport(ctx, start, end)
}

func (uc *ReportUseCase) getReport(ctx context.Context, start, end time.Time) (*models.ReportResponse, error) {
	summary, err := uc.TransactionRepository.GetReportSummary(ctx, start, end)
	if err != nil {
		return nil, &models.Error{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return &models.ReportResponse{
		TotalRevenue:     summary.TotalRevenue,
		TotalTransaction: summary.TotalTransaction,
		BestSellingProduct: models.ReportTopProduct{
			Name:      summary.TopProductName,
			TotalSold: summary.TopProductQty,
		},
	}, nil
}
