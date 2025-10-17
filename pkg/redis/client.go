package redis

import (
	"bufio"
	"net"

	"github.com/decoesp/tamborete/internal/resp"
)

type Client struct {
	conn   net.Conn
	reader *bufio.Reader
}

func New(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}, nil
}

func (c *Client) Set(key, value string) error {
	cmd := resp.Serialize([]interface{}{"SET", key, value})
	_, err := c.conn.Write(cmd)
	return err
}

func (c *Client) Get(key string) (string, error) {
	cmd := resp.Serialize([]interface{}{"GET", key})
	_, err := c.conn.Write(cmd)
	if err != nil {
		return "", err
	}

	parser := resp.NewParser(c.reader)
	response, err := parser.Parse()
	if err != nil {
		return "", err
	}
	return response.(string), nil
}
