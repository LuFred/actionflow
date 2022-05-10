package core

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

//Handler base handler
type Handler struct {
	// context data
	Ctx  *Context
	Data map[interface{}]interface{}
}

//HandlerInterface handler interface
type HandlerInterface interface {
	Init(ctx *Context)
}

//Init 实现HandlerInterface接口
func (h *Handler) Init(ctx *Context) {
	h.Ctx = ctx
	h.Data = make(map[interface{}]interface{})
}

// ServeJSON sends a json response with encoding charset.
func (h *Handler) ServeJSON(encoding ...bool) {
	var (
		hasIndent   = true
		hasEncoding = false
	)
	// if BConfig.RunMode == PROD {
	// 	hasIndent = false
	// }
	if len(encoding) > 0 && encoding[0] {
		hasEncoding = true
	}
	h.Ctx.Response.JSON(h.Data["json"], hasIndent, hasEncoding)
}

// ServeJSONString sends a json string response with encoding charset.
func (h *Handler) ServeJSONString(encoding ...bool) {
	var (
		hasIndent   = true
		hasEncoding = false
	)
	// if BConfig.RunMode == PROD {
	// 	hasIndent = false
	// }
	if len(encoding) > 0 && encoding[0] {
		hasEncoding = true
	}
	h.Ctx.Response.JSONString(h.Data["jsonString"].(string), hasIndent, hasEncoding)
}

// Abort stops controller handler and show the error data if code is defined in ErrorMap or code string.
func (h *Handler) Abort(code string) {
	status, err := strconv.Atoi(code)
	if err != nil {
		status = 200
	}
	h.CustomAbort(status, code)
}

//ErrorAbort 500
func (h *Handler) ErrorAbort(errMsg string) {
	h.Ctx.Response.WriteHeader(http.StatusInternalServerError)
	h.Ctx.Response.Write([]byte(errMsg))
}

//BadRequestAbort 400
func (h *Handler) BadRequestAbort(msg string) {
	//ignore error
	h.Ctx.Response.Header("Content-Type", "application/json")
	h.Ctx.Response.WriteHeader(http.StatusBadRequest)
	em := fmt.Sprintf(`{"message":"%s"}`, msg)
	h.Ctx.Response.Write([]byte(em))

}

//NoContentAbort 204
func (h *Handler) NoContentAbort() {
	h.Ctx.Response.WriteHeader(http.StatusNoContent)
}

//UnauthorizedAbort 401
func (h *Handler) UnauthorizedAbort(msg string) {
	h.Ctx.Response.WriteHeader(http.StatusUnauthorized)
	h.Ctx.Response.Write([]byte(msg))
}

//NotFoundAbort 404
func (h *Handler) NotFoundAbort(msg string) {
	h.Ctx.Response.WriteHeader(http.StatusNotFound)
	h.Ctx.Response.Write([]byte(msg))
}

// CustomAbort stops controller handler and show the error data, it's similar Aborts, but support status code and body.
func (h *Handler) CustomAbort(status int, body string) {
	h.Ctx.Response.WriteHeader(status)
	h.Ctx.Response.Write([]byte(body))

}

//CustomJsonAbort stops controller handler and show the error data,it's similar Aborts,but support status code and json format body.
func (h *Handler) CustomJsonAbort(status int, data interface{}) {
	h.Ctx.Response.Header("Content-Type", "application/json")
	h.Ctx.Response.WriteHeader(status)
	h.Ctx.Response.JSON(data, false, false)
}

//CustomJsonStringAbort stops controller handler and show the error data,it's similar Aborts,but support  json format body.
func (h *Handler) CustomJsonStringAbort(status int, data string) {
	h.Ctx.Response.Header("Content-Type", "application/json")
	h.Ctx.Response.WriteHeader(status)
	h.Ctx.Response.Write([]byte(data))
}

//CustomTextPlainAbort
func (h *Handler) CustomTextPlainAbort(status int, data string) {
	h.Ctx.Response.Header("Content-Type", "text/plain")
	h.Ctx.Response.WriteHeader(status)
	h.Ctx.Response.Write([]byte(data))
}

//Next
func (h *Handler) Next(resp interface{}, err error) {
	if err != nil {
		httpErr, ok := err.(HttpStatusError)
		if ok {
			switch httpErr.Status {
			case http.StatusBadRequest:
				h.CustomJsonStringAbort(httpErr.Status, httpErr.Body)
			default:
				h.CustomAbort(httpErr.Status, httpErr.Body)
			}
		} else {
			h.ErrorAbort(err.Error())
		}
		return

	}
	if resp == nil {
		h.NotFoundAbort("not found")
		return
	}
	if reflect.TypeOf(resp).Kind() == reflect.String {
		switch h.Ctx.Response.ContentType {
		case ResponseContentType_ApplicationJSON:
			h.Data["jsonString"] = resp
			h.ServeJSONString()
		case ResponseContentType_TextPlain:
			h.CustomTextPlainAbort(200, resp.(string))
		default:
			h.Data["jsonString"] = resp
			h.ServeJSONString()
		}

	} else {
		h.Data["json"] = resp
		h.ServeJSON()
	}

}

//VerifyPermission 验证用户操作是否有权限
func (h *Handler) VerifyPermission(operation string) bool {
	var ok bool
	ok = true
	// if h.Ctx.GetUser() == nil || h.Ctx.GetUser().UserID == 0 {
	// 	return ok
	// }
	// url := fmt.Sprintf(conf.ConfigCenterApiServer.GetUserSection, h.Ctx.GetUser().UserOauthIdentity, operation)
	//
	// client := &http.Client{}
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	log.Debug(err)
	// 	log.Errorf("VerifyPermission url:%s err:%v", url, err)
	// 	return ok
	// }
	// //req.Header.Set("Authorization", "Bearer "+token)
	// resp, err := client.Do(req)
	// if err == nil {
	// 	log.Debug(err)
	// 	log.Errorf("VerifyPermission url:%s err:%v", url, err)
	// 	defer resp.Body.Close()
	// }
	// if resp.StatusCode == http.StatusOK {
	// 	body, _ := ioutil.ReadAll(resp.Body)
	// 	if strings.ToLower(string(body[:])) == "true" {
	// 		ok = true
	// 	}
	// }
	return ok
}
