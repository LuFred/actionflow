package appsflow

import (
	"actionflow/worker/appsflow/activities/http"
	"actionflow/worker/appsflow/activities/pageui"
	"actionflow/worker/appsflow/pkg/httputil"
	"encoding/json"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

//TODO Get from the configuration file
var (
	APIGetJobTemplate            = "http://localhost:4000/flows/jobs/%s"
	APIGetJobRunInstanceTemplate = "http://localhost:4000/flows/jobs/%s/runinstance/%s"
	APIGetFlowActionTemplate     = "http://localhost:4000/flows/%s/actions"

	FlowVariableKeyPrefix = "flow.variables."
)

func RunJobRunInstanceWorkflow(ctx workflow.Context, runInstanceId string, jobId string) error {
	workflowID := runInstanceId
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 30 * time.Minute,
	}

	logger := workflow.GetLogger(ctx)
	ctx = workflow.WithActivityOptions(ctx, ao)

	//step1. get job RunInstance
	getJobRunInstanceReq := &httputil.Request{
		Url:  fmt.Sprintf(APIGetJobRunInstanceTemplate, jobId, runInstanceId),
		Verb: httputil.GET,
	}
	getJobRunInstanceResp, err := SendHttp(ctx, getJobRunInstanceReq)
	if err != nil {
		logger.Error("GetJobRunInstance failed.", "error", err.Error())
		return err
	}

	if getJobRunInstanceResp.Err != "" {
		logger.Error("getJobRunInstance failed", "err", getJobRunInstanceResp.Err, "statusCode", getJobRunInstanceResp.StatusCode)
		// TODO: get runInstance with error, update job status to failed
		return nil
	}

	jobRunInstance := JobRunInstance{}
	err = json.Unmarshal(getJobRunInstanceResp.Body, &jobRunInstance)
	if err != nil {
		logger.Error("jobRunInstance  Unmarshal failed", "err", err.Error(), "body", string(getJobRunInstanceResp.Body))
		// TODO: Unmarshal with error, update job status to failed
		return nil
	}

	//step2. get job
	logger.Info(fmt.Sprintf("jobid = %s", jobId))
	getJobReq := &httputil.Request{
		Url:  fmt.Sprintf(APIGetJobTemplate, jobRunInstance.JobId),
		Verb: httputil.GET,
	}
	getJobResp, err := SendHttp(ctx, getJobReq)
	if err != nil {
		logger.Error("GetJob failed.", "error", err.Error())
		return err
	}

	if getJobResp.Err != "" {
		logger.Error("getJobRunInstance failed", "err", getJobResp.Err, "statusCode", getJobResp.StatusCode)
		// TODO: get job with error, update job status to failed
		return nil
	}

	job := &Job{}
	err = json.Unmarshal(getJobResp.Body, job)
	if err != nil {
		logger.Error("job Unmarshal failed", "err", err.Error(), "body", string(getJobResp.Body))
		// TODO: Unmarshal with error, update job status to failed
		return nil
	}

	//step3. get flowActions
	getActionsReq := &httputil.Request{
		Url:  fmt.Sprintf(APIGetFlowActionTemplate, job.FlowId),
		Verb: httputil.GET,
	}

	getActionsResp, err := SendHttp(ctx, getActionsReq)
	if err != nil {
		logger.Error("getActions failed.", "error", err.Error())
		return err
	}

	if getActionsResp.Err != "" {
		logger.Error("getActions failed", "err", getActionsResp.Err, "statusCode", getActionsResp.StatusCode)
		// TODO: get actions with error, update job status to failed
		return nil
	}

	flowAction := &FlowAction{}
	err = json.Unmarshal(getActionsResp.Body, flowAction)
	if err != nil {
		logger.Error("flowAction Unmarshal failed", "err", err.Error(), "body", string(getActionsResp.Body))
		// TODO: Unmarshal with error, update job status to failed
		return nil
	}

	//step 4.2 解析节点动作
	var actions []*ActivityInvocation
	actions = make([]*ActivityInvocation, len(flowAction.Data))
	for i, act := range flowAction.Data {
		actions[i] = &ActivityInvocation{
			Id:          act.Id,
			PerIds:      act.PreIds,
			Name:        act.Name,
			Type:        act.Type,
			DisplayName: act.DisplayName,
		}

		switch ActionType(act.Type) {
		case HTTPAction:
			var cmd http.Command
			err = json.Unmarshal([]byte(act.Command), &cmd)
			if err != nil {
				logger.Error("command Unmarshal failed", "err", err.Error(), "command", act.Command)
				// TODO: Unmarshal with error, update job status to failed
				return nil
			}
			actions[i].Command = &cmd
		case PageAction:
			var cmd pageui.Command
			err = json.Unmarshal([]byte(act.Command), &cmd)
			if err != nil {
				logger.Error("command Unmarshal failed", "err", err.Error(), "command", act.Command)
				// TODO: Unmarshal with error, update job status to failed
				return nil
			}
			actions[i].Command = &cmd
		}
	}

	// step 4.1 解析全部变量
	flowVariable := make(map[string]string)
	if len(jobRunInstance.FlowVariable) > 0 {
		tmpVariables, err := JsonToMap(jobRunInstance.FlowVariable)
		if err != nil {
			logger.Info("JsonToMap failed", "var", jobRunInstance.FlowVariable)
			// TODO: Unmarshal with error, update job status to failed
			return nil
		}

		for k, v := range tmpVariables {
			flowVariable[fmt.Sprintf("%s%s", FlowVariableKeyPrefix, k)] = v
		}
	}

	w := NewWorker(ctx, logger, workflowID, actions, flowVariable)
	err = w.Execute()
	logger.Info(fmt.Sprintf("workflow done"))
	logger.Info(fmt.Sprintf("flow.variables = %v", w.Variables))

	return err
}

func SendHttp(ctx workflow.Context, req *httputil.Request) (*httputil.Response, error) {
	var resp httputil.Response
	var a *http.HttpActivities
	err := workflow.ExecuteActivity(ctx, a.HttpCall, req).Get(ctx, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}