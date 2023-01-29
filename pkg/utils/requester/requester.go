package requester

import "net/http"

type Requester interface {
	Do(*http.Request) (*http.Response, error)
}
