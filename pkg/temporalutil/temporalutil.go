package temporalutil

import (
	"fmt"
	"go.temporal.io/sdk/client"
)

type Client struct {
	client.Client
	HostPort string
}

var (
	tc *Client
)

func InitClient(hostPort string) error {
	var err error
	if tc != nil {
		fmt.Errorf("duplicate initialization")
	}

	tc, err = newClient(hostPort)
	return err
}

func GetClient() (*Client, error) {
	if tc == nil {
		return nil, fmt.Errorf("the client is not initialized")
	}

	return tc, nil
}

func newClient(hostPort string) (*Client, error) {
	cli := &Client{
		HostPort: hostPort,
	}

	c, err := client.NewClient(client.Options{
		HostPort: hostPort,
	})
	if err != nil {
		return nil, err
	}
	cli.Client = c

	return cli, nil
}
