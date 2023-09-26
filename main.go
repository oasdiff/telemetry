package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/gin-gonic/gin"
	"github.com/oasdiff/telemetry/model"
	"github.com/oasdiff/go-common/task"
	"golang.org/x/exp/slog"
)

func main() {

	setupRouter().Run(":8080")
}

func setupRouter() *gin.Engine {

	builder := task.NewTaskBuilder()
	router := gin.Default()
	router.POST(fmt.Sprintf("/%s", model.KeyEvents), func(ctx *gin.Context) {

		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(ctx.Request.Body); err != nil {
			slog.Error("failed to read request body", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer ctx.Request.Body.Close()

		body := buf.Bytes()
		size := len(body)
		if size == 0 {
			slog.Info("empty body")
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if size > 102400 { // 100KB
			slog.Info("client sent payload > 100KB", "size", size)
			ctx.Writer.WriteHeader(http.StatusBadRequest)
			return
		}

		client, err := cloudtasks.NewClient(context.Background())
		if err != nil {
			slog.Error("failed to create cloud task client", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer client.Close()
		_, err = client.CreateTask(ctx, builder.CreateRequest(body))
		if err != nil {
			slog.Error("failed to create cloud task", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx.Writer.WriteHeader(http.StatusCreated)
	})

	return router
}
