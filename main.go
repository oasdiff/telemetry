package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oasdiff/telemetry/model"
	"golang.org/x/exp/slog"
)

func main() {

	setupRouter().Run(":8080")
}

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.POST("/commands", func(ctx *gin.Context) {
		var t model.Telemetry
		if err := json.NewDecoder(ctx.Request.Body).Decode(&t); err != nil {
			err = fmt.Errorf("failed to decode create telemetry request body with '%v'", err)
			slog.Info(err.Error())
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		payload, err := json.MarshalIndent(t, "", "    ")
		if err != nil {
			slog.Error("failed to 'MarshalIndent'", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info(string(payload))

		ctx.Writer.WriteHeader(http.StatusCreated)
	})

	return router
}
