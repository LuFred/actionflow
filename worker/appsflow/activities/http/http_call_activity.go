package http

import (
	"actionflow/worker/appsflow/pkg/httputil"
	"context"
	"fmt"
	"log"
	"time"
)

type HttpActivities struct {}

type Command struct {
	Url               string       `json:"url"`
	Verb              string       `json:"verb"`
	SuccessStatusCode int          `json:"successStatusCode"`
	Parameters        []*Parameter `json:"parameters"`
}

type Parameter struct {
	Name      string `json:"name"`
	DataType  string `json:"dataType" description:"text | integer | boolean"`
	ReceiveIn string `json:"receiveIn" description:"body | query | header"`
	ValueType string `json:"valueType" description:"reference"`
	Value     string `json:"value"`
}

func (a *HttpActivities) HttpCall(ctx context.Context, req httputil.Request) (*httputil.Response, error) {
	for _, v := range req.Parameters {
		log.Printf("[HttpCall] par =%+v", v)
	}
	client := &httputil.HttpClient{}
	var resp *httputil.Response
	for {
		select {
		case <-ctx.Done():
			log.Println("activity Done")
			return nil, fmt.Errorf("activity canceled")
		case <-time.After(time.Second * 5):
		}

		bodyBytes, statusCode, err := client.Send(&req)

		if err != nil {
			resp.Err = err.Error()
		}
		log.Printf("[HttpCall] request = %+v status =%d", req, statusCode)
		if statusCode == 200 || statusCode == 201 {
			resp = &httputil.Response{
				StatusCode: statusCode,
				Body:       bodyBytes,
			}
			break
		}
	}
	return resp, nil
}
