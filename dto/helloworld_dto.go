package dto

type HelloworldRequest struct {
	Name string `json:"name" transport:"query" required:"true" description:"name"`
}

type HelloworldResponse struct {
	Message string `json:"message" description:"helloworld messge"`
}
