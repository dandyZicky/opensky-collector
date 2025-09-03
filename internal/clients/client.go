// Package clients provides interface for HTTP operations
package clients

import "net/http"

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
