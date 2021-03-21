package commands

import "net/http"

type Http interface {
	Get(url string) (*http.Response, error)
}
