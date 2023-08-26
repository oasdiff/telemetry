package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oasdiff/go-common/env"
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

		payload, err := json.Marshal(events)
		if err != nil {
			slog.Error("failed to 'Marshal'", "events", len(events), "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		slog.Info(string(payload))

		ctx.Writer.WriteHeader(http.StatusCreated)
	})
	router.GET("/test", func(ctx *gin.Context) {
		client, err := cloudtasks.NewClient(context.Background())
		if err != nil {
			slog.Error("failed to create cloud task client", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer client.Close()

		request := &taskspb.CreateTaskRequest{
			Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetGCPProject(), env.GetGCPLocation(), env.GetGCPQueue()),
			Task: &taskspb.Task{
				MessageType: &taskspb.Task_HttpRequest{
					HttpRequest: &taskspb.HttpRequest{
						HttpMethod: taskspb.HttpMethod_POST,
						Url:        env.GetTaskSubscriberUrl(),
					},
				},
			},
		}
		request.Task.GetHttpRequest().Body = streamToByte(ctx.Request.Body)

		_, err = client.CreateTask(ctx, request)
		if err != nil {
			slog.Error("failed to create cloud task", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	return router
}

func streamToByte(stream io.Reader) []byte {

	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}
