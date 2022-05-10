package appsflow

import (
	"actionflow/worker/appsflow/activities/http"
	"actionflow/worker/appsflow/activities/pageui"
	"actionflow/worker/appsflow/pkg/httputil"
	"fmt"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
	"sync"
)

type Worker struct {
	id string
	runc      workflow.Channel
	completec workflow.Channel
	mu        *sync.RWMutex

	err   error
	errMu *sync.RWMutex

	Variables      map[string]string
	actions        []*ActivityInvocation
	actionMap      map[string]*ActivityInvocation
	actionCount    int
	maxParallelNum int // Maximum number of parallel tasks. default 10

	completeNum   int
	completeNumMu *sync.RWMutex

	curParallelNum   int
	curParallelNumMu *sync.RWMutex

	ctx      workflow.Context
	cancel   workflow.CancelFunc
	canceled bool

	lg log.Logger

	edges map[string][]string
	indeg map[string]int
}

type ActivityInvocation struct {
	Id               string
	PerIds           []string
	Name             string
	Type             string
	DisplayName      string
	Command          interface{}
	Condition        string
	SuccessCondition string
	Error            error
}

var WorkflowActionMaxParallelNum = 5

type ActionType string

const (
	HTTPAction ActionType = "HTTP"
	PageAction ActionType = "Page"
)

type ParameterValueType string

const (
	ValueTypeReference ParameterValueType = "reference"
	ValueTypeConstant  ParameterValueType = "constant"
)

func NewWorker(ctx workflow.Context, lg log.Logger,workflowID string, actions []*ActivityInvocation, variables map[string]string) *Worker {
	childCtx, cancel := workflow.WithCancel(ctx)
	w := &Worker{
		id: workflowID,
		runc:             workflow.NewChannel(childCtx),
		completec:        workflow.NewChannel(childCtx),
		actions:          actions,
		actionMap:        make(map[string]*ActivityInvocation),
		Variables:        variables,
		actionCount:      len(actions),
		completeNum:      0,
		maxParallelNum:   WorkflowActionMaxParallelNum,
		curParallelNumMu: new(sync.RWMutex),
		completeNumMu:    new(sync.RWMutex),
		err:              nil,
		errMu:            new(sync.RWMutex),
		ctx:              childCtx,
		cancel:           cancel,
		canceled:         false,
		mu:               new(sync.RWMutex),
		lg:               lg,
		edges:            make(map[string][]string, len(actions)),
		indeg:            make(map[string]int, len(actions)),
	}

	for _, info := range w.actions {
		w.actionMap[info.Id] = info
		w.indeg[info.Id] = 0
		for _, pid := range info.PerIds {
			w.edges[pid] = append(w.edges[pid], info.Id)
			w.indeg[info.Id]++
		}
	}

	workflow.Go(w.ctx, func(ctx workflow.Context) {
		w.run(ctx)
	})

	return w
}

// setError Only log the first activity error
func (w *Worker) setError(err error) {
	w.errMu.Lock()
	defer w.errMu.Unlock()
	if w.err == nil {
		w.err = err
		w.Cancel()
	}
}

func (w *Worker) Cancel() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.canceled {
		w.cancel()
		w.canceled = true
	}
}

func (w *Worker) Execute() error {
	defer w.runc.Close()
	var q []string
	for k, v := range w.indeg {
		if v == 0 {
			q = append(q, k)
		}
	}
	for _, v := range q {
		w.runc.Send(w.ctx, v)
	}

	selector := workflow.NewSelector(w.ctx)
	selector.AddReceive(w.ctx.Done(), func(c workflow.ReceiveChannel, more bool) {
		// workflow cancel
		w.lg.Info("workflow done")
	})
	selector.AddReceive(w.completec, func(c workflow.ReceiveChannel, more bool) {
		if !more {
			w.lg.Info("completec closed")
			return
		}

		var completeId string
		c.Receive(w.ctx, &completeId)
		w.completeNumMu.Lock()
		defer w.completeNumMu.Unlock()
		w.completeNum++
		for _, v := range w.edges[completeId] {
			w.indeg[v]--
			if w.indeg[v] == 0 {
				w.runc.Send(w.ctx, v)
			}
		}
	})

	for w.completeNum < w.actionCount && w.err == nil && !w.canceled {
		selector.Select(w.ctx)
	}

	return nil
}

func (w *Worker) run(ctx workflow.Context) {
	for w.err == nil {
		var id string
		if closed := w.runc.Receive(ctx, &id); !closed {
			return
		}

		workflow.Go(ctx, func(ctx workflow.Context) {
			w.curParallelNumMu.Lock()
			w.curParallelNum++
			w.curParallelNumMu.Unlock()
			w.execute(ctx, id)
		})
	}
}

func (w *Worker) execute(ctx workflow.Context, id string) {
	//TODO set ActivityOptions
	//ao := workflow.ActivityOptions{
	//	StartToCloseTimeout: 10 * time.Second,
	//	RetryPolicy: &temporal.RetryPolicy{
	//		MaximumAttempts: 1,
	//	},
	//}
	//ctx = workflow.WithActivityOptions(ctx, ao)
	var err error
	switch ActionType(w.actionMap[id].Type) {
	case HTTPAction:
		err = w.executeHttpActivity(ctx, w.actionMap[id])
	case PageAction:
		err = w.executePageUIActivity(ctx, w.actionMap[id])
	default:
		// TODO: err Incorrect action type
	}
	// TODO update activity status
	if err != nil {
		w.actionMap[id].Error = err
		w.setError(err)
	}

	w.curParallelNumMu.Lock()
	w.curParallelNum--
	w.curParallelNumMu.Unlock()

	w.completec.Send(ctx, id)
	return
}

func (w *Worker) executeHttpActivity(ctx workflow.Context, ai *ActivityInvocation) error {
	var err error

	if dc, ok := ai.Command.(*http.Command); ok {
		req := &httputil.Request{
			Url:        dc.Url,
			Verb:       httputil.HTTPVerb(dc.Verb),
			Parameters: make([]*httputil.HttpParameter, len(dc.Parameters)),
		}

		// parameter conversion
		for i, p := range dc.Parameters {
			req.Parameters[i] = &httputil.HttpParameter{
				Name:      p.Name,
				ReceiveIn: httputil.HTTPReceivein(p.ReceiveIn),
				DataType:  p.DataType,
			}
			switch ParameterValueType(p.ValueType) {
			case ValueTypeReference:
				req.Parameters[i].Value = w.Variables[p.Value]
			case ValueTypeConstant:
				req.Parameters[i].Value = p.Value
			default:
				err = NewWorkflowError(fmt.Sprintf("Incorrect parameter type(%s)", p.ValueType), ERROR_TYPE_ACTIVITY)
				return err
			}
		}

		var a *http.HttpActivities
		f := workflow.ExecuteActivity(ctx, a.HttpCall, req)

		var resp httputil.Response
		err = f.Get(ctx, &resp)
		if err != nil {
			if err == workflow.ErrCanceled {
				err = NewWorkflowError(err.Error(), ERROR_TYPE_CANCELED)
			} else {
				err = NewWorkflowError(err.Error(), ERROR_TYPE_ACTIVITY)
			}
		}

		w.Variables[fmt.Sprintf("node.%s.response", ai.Name)] = string(resp.Body)
	} else {
		err = NewWorkflowError("command Unmarshal failed", ERROR_TYPE_ACTIVITY)
	}

	return err
}

func (w *Worker) executePageUIActivity(ctx workflow.Context, ai *ActivityInvocation) error {
	var err error
	fmt.Println(fmt.Sprintf("start run ============= = %+v", ai))
	if cmd, ok := ai.Command.(*pageui.Command); ok {
		fmt.Println(fmt.Sprintf("start run ============= = %+v", cmd))
			switch ParameterValueType(cmd.ValueType) {
			case ValueTypeReference:
				cmd.Value = w.Variables[cmd.Value]
			case ValueTypeConstant:
			default:
				err = NewWorkflowError(fmt.Sprintf("Incorrect parameter type(%s)", cmd.ValueType), ERROR_TYPE_ACTIVITY)
				return err
			}

		var a *pageui.PageUIActivities
		f := workflow.ExecuteActivity(ctx, a.ActivityWaitCallback, w.id, cmd.Value)

		var resp string
		err = f.Get(ctx, &resp)
		if err != nil {
			if err == workflow.ErrCanceled {
				err = NewWorkflowError(err.Error(), ERROR_TYPE_CANCELED)
			} else {
				err = NewWorkflowError(err.Error(), ERROR_TYPE_ACTIVITY)
			}
		}
		fmt.Println(fmt.Sprintf("executePageUIActivity = %s", resp))
		w.Variables[fmt.Sprintf("node.%s.response", ai.Name)] = resp
	} else {
		err = NewWorkflowError("command Unmarshal failed", ERROR_TYPE_ACTIVITY)
	}

	return err
}
