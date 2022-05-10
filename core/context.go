package core

import (
	"net/http"
)

//Context b
type Context struct {
	Request   *Request
	Response  *Response
	OauthInfo *OauthInfo
	SrvCtx    *Server
}

//GetUser 获取user对象
func (c *Context) GetOauthInfo() *OauthInfo {
	if c.OauthInfo == nil {
		return nil
	}
	return c.OauthInfo
}

//GetUserID 获取userid
func (c *Context) GetUserID() string {
	if c.OauthInfo == nil {
		return ""
	}
	return c.OauthInfo.UserId
}

//NewContext 创建context
func NewContext(r *http.Request, w http.ResponseWriter, s *Server) *Context {
	c := new(Context)
	c.Request = &Request{
		Request: r,
	}
	c.Response = &Response{
		ResponseWriter: w,
	}
	c.SrvCtx = s

	return c
}

//InitUser 初始化user对象
func (c *Context) setOauthInfo(u *OauthInfo) {
	if u != nil {
		c.OauthInfo = u
	}
}
