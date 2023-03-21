package reports_used_car

import (
	"context"
	"encoding/json"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	"time"
)

type SubreportResponseOwnershipHistory struct {
	VIN    string
	Owners DatedList
	SubreportResponse
}

type SubreportOwnershipHistoryHandler struct {
	RabbitClient *rabbitclient.RabbitClient
}

func (h *SubreportOwnershipHistoryHandler) Handle(ctx context.Context, data []byte) error {
	report := &Report{}
	err := json.Unmarshal(data, report)
	if err != nil {
		return err
	}
	ownershipHistory, _ := OwnershipHistoryByVINs[report.VIN]
	response := SubreportResponseOwnershipHistory{
		VIN:    report.VIN,
		Owners: ownershipHistory,
		SubreportResponse: SubreportResponse{
			ReportID: report.ID,
			Type:     "ownership-history",
			Date:     time.Now(),
		},
	}
	return h.RabbitClient.PublishMessage(ctx, "reports.used-car", response)
}
