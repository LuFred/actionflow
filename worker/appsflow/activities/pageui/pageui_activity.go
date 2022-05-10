package pageui

import (
	"context"
	"errors"
	"fmt"
	"go.temporal.io/sdk/activity"
	"io/ioutil"
	"net/http"
	"net/url"
)

type PageUIActivities struct{}

type Command struct {
	ValueType string `json:"valueType" description:"reference"`
	Value     string `json:"value"`
}

func (a *PageUIActivities) ActivityWaitCallback(ctx context.Context, workflowID string, sessionID string) (string, error) {
	if len(sessionID) == 0 {
		return "", errors.New("session id is empty")
	}
	logger := activity.GetLogger(ctx)
	logger.Info("[ActivityWaitCallback]", "workflowID", workflowID, "sessionID", sessionID)

	activityInfo := activity.GetInfo(ctx)
	formData := url.Values{}
	formData.Add("task_token", string(activityInfo.TaskToken))
	fmt.Println(fmt.Sprintf("id =%s.%s   token = %s", workflowID, sessionID, string(activityInfo.TaskToken)))
	registerCallbackURL := "http://localhost:8099/registerCallback?id=" + fmt.Sprintf("%s.%s", workflowID, sessionID)
	resp, err := http.PostForm(registerCallbackURL, formData)
	if err != nil {
		logger.Info("ActivityWaitCallback failed to register callback.", "Error", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return "", err
	}
	status := string(body)
	if status == "SUCCEED" {
		// register callback succeed
		logger.Info("Successfully registered callback.", "sessionID", sessionID)

		// ErrActivityResultPending is returned from activity's execution to indicate the activity is not completed when it returns.
		// activity will be completed asynchronously when Client.CompleteActivity() is called.
		return "", activity.ErrResultPending
	}

	logger.Warn("Register callback failed.", "ExpenseStatus", status)
	return "", fmt.Errorf("register callback failed status:%s", status)
}
