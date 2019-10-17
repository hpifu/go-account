package account

import (
	"time"

	"github.com/hpifu/go-kit/hhttp"
)

type Client struct {
	client  *hhttp.HttpClient
	address string
}

func NewClient(address string, maxConn int, connTimeout time.Duration, recvTimeout time.Duration) *Client {
	return &Client{
		address: address,
		client:  hhttp.NewHttpClient(maxConn, connTimeout, recvTimeout),
	}
}
