package reports_used_car

import (
	"context"
	"encoding/json"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	"time"
)

type SubreportResponseAccidents struct {
	VIN       string
	Accidents DatedList
	SubreportResponse
}

type SubreportAccidentsHandler struct {
	RabbitClient *rabbitclient.RabbitClient
}

func (h *SubreportAccidentsHandler) Handle(ctx context.Context, data []byte) error {
	report := &Report{}
	err := json.Unmarshal(data, report)
	if err != nil {
		return err
	}
	accidents, _ := AccidentsByVINs[report.VIN]
	response := SubreportResponseAccidents{
		VIN:       report.VIN,
		Accidents: accidents,
		SubreportResponse: SubreportResponse{
			ReportID: report.ID,
			Type:     "accidents",
			Date:     time.Now(),
		},
	}
	return h.RabbitClient.PublishMessage(ctx, "reports.used-car", response)
}
