package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

//Response Response扩展
type Response struct {
	http.ResponseWriter
	Status      int
	Started     bool
	EnableGzip  bool
	ContentType ResponseContentType
}
type ResponseContentType int

const (
	ResponseContentType_ApplicationJSON ResponseContentType = iota // value --> 0
	ResponseContentType_TextPlain
)

// Header set response header
func (w *Response) Header(key, val string) error {
	w.ResponseWriter.Header().Set(key, val)
	return nil

}

// WriteBody reads  writes content to writer by the specific encoding(gzip/deflate)
func WriteBody(writer io.Writer, content []byte) error {
	_, err := writer.Write(content)
	return err

}

// Body sets response body content.
// if EnableGzip, compress content string.
// it sends out response body directly.
func (w *Response) Body(content []byte) error {

	var buf = &bytes.Buffer{}
	//todo gzip compress

	w.Header("Content-Length", strconv.Itoa(len(content)))
	// Write status code if it has been set manually
	// Set it to 0 afterwards to prevent "multiple response.WriteHeader calls"
	if w.Status != 0 {
		w.ResponseWriter.WriteHeader(w.Status)
		w.Status = 0
	} else {
		w.Started = true
	}

	WriteBody(buf, content)
	io.Copy(w.ResponseWriter, buf)
	return nil
}

//JSONString Send json string format data to the client
func (w *Response) JSONString(data string, hasIndent bool, coding bool) error {

	w.Header("Content-Type", "application/json")
	var content []byte
	if coding {
		content = []byte(stringsToJSON(data))
	} else {
		content = []byte(data)
	}
	return w.Body(content)
}

//JSON Send json format data to the client
func (w *Response) JSON(data interface{}, hasIndent bool, coding bool) error {

	w.Header("Content-Type", "application/json")
	var content []byte
	var err error
	if hasIndent {
		content, err = json.MarshalIndent(data, "", " ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		http.Error(w.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return err
	}
	if coding {
		content = []byte(stringsToJSON(string(content)))
	}
	return w.Body(content)
}
func stringsToJSON(str string) string {
	var jsons bytes.Buffer
	for _, r := range str {
		rint := int(r)
		if rint < 128 {
			jsons.WriteRune(r)
		} else {
			jsons.WriteString("\\u")
			jsons.WriteString(strconv.FormatInt(int64(rint), 16))
		}
	}
	return jsons.String()
}
