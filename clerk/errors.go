package clerk

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Response *http.Response
	Errors   []Error `json:"errors"`
}

type Error struct {
	Message     string      `json:"message"`
	LongMessage string      `json:"long_message"`
	Code        string      `json:"code"`
	Meta        interface{} `json:"meta,omitempty"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		e.Response.Request.Method, e.Response.Request.URL,
		e.Response.StatusCode, e.Errors)
}
