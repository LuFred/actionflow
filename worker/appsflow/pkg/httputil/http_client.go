package httputil

import (
	"fmt"
	"github.com/ddliu/go-httpclient"
)

type HTTPVerb string

const (
	GET    HTTPVerb = "get"
	POST   HTTPVerb = "post"
	PUT    HTTPVerb = "put"
	DELETE HTTPVerb = "delete"
	PATCH  HTTPVerb = "patch"
)

type HTTPReceivein string

const (
	ReceiveinBody   HTTPReceivein = "body"
	ReceiveinQuqery HTTPReceivein = "query"
	ReceiveinHeader HTTPReceivein = "header"
)

type Response struct {
	StatusCode int
	Body       []byte
	Err        string
}

type Request struct {
	Parameters []*HttpParameter `json:"parameters"`
	Verb       HTTPVerb         `json:"verb"`
	Url        string           `json:"url"`
}

type HttpParameter struct {
	Name      string        `json:"name"`
	DataType  string        `json:"dataType" description:"text | integer | boolean"`
	ReceiveIn HTTPReceivein `json:"receiveIn" description:"body | query | header"`
	Value     string        `json:"value"`
}


type HttpClient struct{}

func (c *HttpClient) Send(req *Request) ([]byte, int, error) {
	var err error
	query := map[string]string{}
	header := map[string]string{}
	body := map[string]string{}
	if len(req.Parameters) > 0 {
		for _, par := range req.Parameters {
			switch par.ReceiveIn {
			case ReceiveinHeader:
				header[par.Name] = par.Value
			case ReceiveinQuqery:
				query[par.Name] = par.Value
			case ReceiveinBody:
				body[par.Name] = par.Value
			default:
				return nil, 400, fmt.Errorf("unsupported request parameter: name(%s) ReceiveIn(%s)", par.Name, par.ReceiveIn)
			}
			query[par.Name] = par.Value
		}
	}

	var reply *httpclient.Response
	hc := httpclient.WithHeaders(header)
	switch req.Verb {
	case DELETE:
		reply, err = hc.Delete(req.Url, query)
	case PUT:
		reply, err = hc.PutJson(req.Url, query)
	case PATCH:
		reply, err = hc.PatchJson(req.Url, query)
	case POST:
		reply, err = hc.PostJson(req.Url, query)
	case GET:
		reply, err = hc.Get(req.Url, query)
	default:
		return nil, 500, fmt.Errorf("unsupported request type")
	}

	if err != nil {
		return nil, 500, err
	}

	bodyBytes, err := reply.ReadAll()
	defer reply.Body.Close()

	return bodyBytes, reply.StatusCode, nil
}
