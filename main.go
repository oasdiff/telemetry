package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oasdiff/telemetry/model"
	"golang.org/x/exp/slog"
)

func main() {

	setupRouter().Run(":8080")
}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.POST(fmt.Sprintf("/%s", model.KeyEvents), func(ctx *gin.Context) {

		var events map[string][]*model.Telemetry
		if err := ctx.BindJSON(&events); err != nil {
			return
		}

		if len(events) > 1 {
			slog.Info("user sent more than 1 telemetry events", "count", len(events))
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		for _, currEvent := range events[model.KeyEvents] {
			currEvent.Id = uuid.NewString()
		}

		payload, err := json.MarshalIndent(events, "", "    ")
		if err != nil {
			slog.Error("failed to 'MarshalIndent'", "events", len(events), "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info(string(payload))

		ctx.Writer.WriteHeader(http.StatusCreated)
	})

	return router
}
