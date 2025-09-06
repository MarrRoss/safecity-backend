package http

import (
	"github.com/go-resty/resty/v2"
	"net"
	"net/http"
	"time"
)

type Client struct {
	*resty.Client
}

func NewClient(baseURL string) *Client {
	client := resty.New()
	client.
		SetBaseURL(baseURL).
		SetRetryCount(5).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(3 * time.Second).
		AddRetryCondition(
			func(r *resty.Response, err error) bool {
				if err != nil {
					return true
				}
				switch r.StatusCode() {
				case http.StatusServiceUnavailable:
					return true
				case http.StatusGatewayTimeout:
					return true
				case http.StatusInternalServerError:
					return true
				case http.StatusTooManyRequests:
					return true
				default:
					return false
				}
			},
		)

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	client.SetTransport(&http.Transport{
		DialContext:           dialer.DialContext,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxConnsPerHost:       60,
		MaxIdleConnsPerHost:   60,
		MaxIdleConns:          60,
	})
	return &Client{Client: client}
}
