package reports_used_car

import (
	"context"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kiaplayer/rabbitmq-report-system/internal/common"
	rabbitclient "github.com/kiaplayer/rabbitmq-report-system/internal/rabbit-client"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

func SetupRoutes(r *gin.Engine, rabbitClient *rabbitclient.RabbitClient) {
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	r.SetFuncMap(template.FuncMap{
		"add": func(a int, b int) int {
			return a + b
		},
	})

	r.LoadHTMLGlob("internal/reports/used-car/templates/*")
	r.Static("/static", "internal/reports/used-car/static")

	r.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		lastResult := session.Get("LastResult")
		session.Delete("LastResult")
		_ = session.Save()
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"pageTitle":  "Отчеты о подержанных машинах",
			"reports":    common.SortMapByKeys(Reports),
			"lastResult": lastResult,
		})
	})

	r.POST("/", func(c *gin.Context) {
		vin := c.PostForm("VIN")
		reportID := strconv.FormatInt(time.Now().Unix(), 10) + uuid.New().String()
		Reports[reportID] = Report{
			ID:   reportID,
			Date: time.Now(),
			VIN:  vin,
		}
		err := rabbitClient.PublishMessage(context.Background(), "reports.used-car", Reports[reportID])
		var lastResult string
		if err != nil {
			lastResult = "При построении отчета возникла ошибка: " + err.Error()
		} else {
			lastResult = "Отчет будет сформирован в ближайшее время"
		}
		session := sessions.Default(c)
		session.Set("LastResult", lastResult)
		_ = session.Save()
		c.Redirect(302, "/")
	})
}
