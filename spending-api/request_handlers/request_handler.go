package request_handlers

import "net/http"

type RequestHandler interface {
	Handle(writer http.ResponseWriter, request *http.Request)
}
