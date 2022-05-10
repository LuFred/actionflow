package dto

type ResponseError struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}
