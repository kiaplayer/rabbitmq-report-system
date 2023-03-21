package reports_used_car

import (
	"context"
	"encoding/json"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	"time"
)

type SubreportResponseLimitsWantedInfo struct {
	VIN              string
	LimitsWantedInfo string
	SubreportResponse
}

type SubreportLimitsWantedInfoHandler struct {
	RabbitClient *rabbitclient.RabbitClient
}

func (h *SubreportLimitsWantedInfoHandler) Handle(ctx context.Context, data []byte) error {
	report := &Report{}
	err := json.Unmarshal(data, report)
	if err != nil {
		return err
	}
	limitsWantedInfo, _ := LimitsWantedInfoByVINs[report.VIN]
	response := SubreportResponseLimitsWantedInfo{
		VIN:              report.VIN,
		LimitsWantedInfo: limitsWantedInfo,
		SubreportResponse: SubreportResponse{
			ReportID: report.ID,
			Type:     "limits-wanted-info",
			Date:     time.Now(),
		},
	}
	return h.RabbitClient.PublishMessage(ctx, "reports.used-car", response)
}
