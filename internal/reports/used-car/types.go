package reports_used_car

import (
	"context"
	"encoding/json"
	"time"
)

type DatedList []struct {
	Date string
	Item string
}

type Report struct {
	ID                          string
	Date                        time.Time
	VIN                         string
	OwnershipHistory            DatedList
	Accidents                   DatedList
	LimitsWantedInfo            string
	OwnershipHistoryWasReceived bool
	AccidentsWasReceived        bool
	LimitsWantedInfoWasReceived bool
}

type SubreportResponse struct {
	ReportID string
	Type     string
	Date     time.Time
}

type ReportResultsHandler struct{}

func (h *ReportResultsHandler) Handle(_ context.Context, data []byte) error {
	taskResponse := &SubreportResponse{}
	err := json.Unmarshal(data, taskResponse)
	if err != nil {
		return err
	}
	if report, found := Reports[taskResponse.ReportID]; found {
		switch taskResponse.Type {
		case "accidents":
			taskResponse := &SubreportResponseAccidents{}
			err = json.Unmarshal(data, taskResponse)
			if err != nil {
				return err
			}
			report.Accidents = taskResponse.Accidents
			report.AccidentsWasReceived = true
			Reports[taskResponse.ReportID] = report
		case "limits-wanted-info":
			taskResponse := &SubreportResponseLimitsWantedInfo{}
			err = json.Unmarshal(data, taskResponse)
			if err != nil {
				return err
			}
			report.LimitsWantedInfo = taskResponse.LimitsWantedInfo
			report.LimitsWantedInfoWasReceived = true
			Reports[taskResponse.ReportID] = report
		case "ownership-history":
			taskResponse := &SubreportResponseOwnershipHistory{}
			err = json.Unmarshal(data, taskResponse)
			if err != nil {
				return err
			}
			report.OwnershipHistory = taskResponse.Owners
			report.OwnershipHistoryWasReceived = true
			Reports[taskResponse.ReportID] = report
		}
	}
	return nil
}
