package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	taskspb "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/gin-gonic/gin"
	"github.com/oasdiff/go-common/env"
	"github.com/oasdiff/telemetry/model"
	"golang.org/x/exp/slog"
)

func main() {

	setupRouter().Run(":8080")
}

func setupRouter() *gin.Engine {

	builder := newTaskBuilder()
	router := gin.Default()
	router.POST(fmt.Sprintf("/%s", model.KeyEvents), func(ctx *gin.Context) {
		client, err := cloudtasks.NewClient(context.Background())
		if err != nil {
			slog.Error("failed to create cloud task client", "error", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer client.Close()

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

type taskBuilder struct {
	requestParent string
	subscriberUrl string
}

func newTaskBuilder() *taskBuilder {

	return &taskBuilder{
		requestParent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", env.GetGCPProject(), env.GetGCPLocation(), env.GetGCPQueue()),
		subscriberUrl: env.GetTaskSubscriberUrl(),
	}
}

func (b *taskBuilder) CreateRequest(body []byte) *taskspb.CreateTaskRequest {

	return &taskspb.CreateTaskRequest{
		Parent: b.requestParent,
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					HttpMethod: taskspb.HttpMethod_POST,
					Url:        b.subscriberUrl,
					Body:       body,
				},
			},
		},
	}
}
