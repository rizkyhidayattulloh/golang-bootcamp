package http

import (
	"kasir-api/internal/models"
	"kasir-api/internal/usecase"
	"kasir-api/internal/util"
	"net/http"
)

type ReportController struct {
	ReportUseCase *usecase.ReportUseCase
}

func NewReportController(reportUseCase *usecase.ReportUseCase) *ReportController {
	return &ReportController{
		ReportUseCase: reportUseCase,
	}
}

func (rc *ReportController) Today(w http.ResponseWriter, r *http.Request) error {
	report, err := rc.ReportUseCase.Today(r.Context())
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.ReportResponse]{
		Data: report,
	})

	return nil
}

func (rc *ReportController) ByDateRange(w http.ResponseWriter, r *http.Request) error {
	query := r.URL.Query()
	startDate := query.Get("start_date")
	endDate := query.Get("end_date")

	report, err := rc.ReportUseCase.ByDateRange(r.Context(), startDate, endDate)
	if err != nil {
		return err
	}

	util.EncodeJSON(w, http.StatusOK, models.WebResponse[*models.ReportResponse]{
		Data: report,
	})

	return nil
}
