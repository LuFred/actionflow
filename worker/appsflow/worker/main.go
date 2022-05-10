package main

import (
	"actionflow/pkg/logutil"
	"actionflow/worker/appsflow"
	"actionflow/worker/appsflow/activities/http"
	"actionflow/worker/appsflow/activities/pageui"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

var lg *zap.Logger

const TaskQueue = "actionFlow-job-run-instance"

func init() {
	lg = logutil.NewLogger(logutil.LoggerConfig{
		Service:     "job-runInstance-worker",
		Loglevel:    "debug",
		WriteType:   logutil.StdoutAndFile,
		EncoderType: logutil.Json,
	})
}

func main() {
	c, err := client.NewClient(client.Options{
		HostPort: "52.130.153.91:7233",
	})

	if err != nil {
		lg.Panic("Unable to create client")
	}
	defer c.Close()

	w := worker.New(c, TaskQueue, worker.Options{})

	w.RegisterWorkflow(appsflow.RunJobRunInstanceWorkflow)
	w.RegisterActivity(&http.HttpActivities{})
	w.RegisterActivity(&pageui.PageUIActivities{})

	lg.Info("started worker")

	err = w.Run(worker.InterruptCh())
	if err != nil {
		lg.Error(err.Error())
	}
}
