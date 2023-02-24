package jsonrpc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	DefaultTimeoutInSecond = 30
)

type Client struct {
	conn net.Conn

	idCounter uint32

	encoder *json.Encoder
	decoder *json.Decoder

	msgs          chan *msgWrapper
	responseChans map[uint32]chan interface{}
	resps         chan map[string]interface{}
}

type msgWrapper struct {
	method   string
	params   interface{}
	response chan interface{}
}

func NewClient(conn net.Conn) *Client {
	e := json.NewEncoder(conn)
	e.SetIndent("", "\t")

	return &Client{
		conn: conn,

		// Get lower 5 digits of the current process ID as the prefix of the request ID
		idCounter: uint32((os.Getpid() % (1 << 6)) * 10000),

		encoder: e,
		decoder: json.NewDecoder(bufio.NewReader(conn)),

		msgs:          make(chan *msgWrapper, 1024),
		responseChans: make(map[uint32](chan interface{})),
		resps:         make(chan map[string]interface{}, 1024),
	}
}

func (c *Client) SendMsgWithTimeout(method string, params interface{}, timeoutInSec int) (res interface{}, err error) {
	id := atomic.AddUint32(&c.idCounter, 1)
	req := NewRequest(id, method, params)
	var resp Response

	defer func() {
		// TODO: Don't know why https://go.dev/play/p/bZ860T6CtcU
		if resp.ErrorInfo == nil {
			err = nil
		}
		err = errors.Wrapf(err, "error sending message, id %v, method %v, params %+v", id, method, params)

		// For debug purpose
		stdenc := json.NewEncoder(os.Stdin)
		stdenc.SetIndent("", "\t")
		stdenc.Encode(req)
		stdenc.Encode(&resp)
	}()

	if err = c.encoder.Encode(req); err != nil {
		return nil, err
	}

	for count := 0; count <= timeoutInSec; count++ {
		if c.decoder.More() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err = c.decoder.Decode(&resp); err != nil {
		return nil, err
	}
	if resp.Id != id {
		return nil, fmt.Errorf("received response but the id mismatching, request id %d, response id %d", id, resp.Id)
	}

	return resp.Result, resp.ErrorInfo
}

func (c *Client) SendMsg(method string, params interface{}) (interface{}, error) {
	return c.SendMsgWithTimeout(method, params, DefaultTimeoutInSecond)
}

func (c *Client) SendCommand(method string, params interface{}) (interface{}, error) {
	return c.SendMsg(method, params)
}

func (c *Client) InitAsync() chan error {
	errChan := make(chan error, 5)
	go c.dispatcher()
	go c.read(errChan)
	return errChan
}

func (c *Client) SendMsgAsync(method string, params interface{}, responseChan chan interface{}) {
	msg := &msgWrapper{method: method,
		params:   params,
		response: responseChan}

	c.msgs <- msg
}

func (c *Client) handleSend(msg *msgWrapper) {
	id := atomic.AddUint32(&c.idCounter, 1)

	req := NewRequest(id, msg.method, msg.params)
	c.responseChans[id] = msg.response

	if err := c.encoder.Encode(req); err != nil {
		logrus.Errorf("error encoding during handleSend: %v", err)
	}
}

func (c *Client) handleRecv(obj map[string]interface{}) {
	fid, ok := obj["id"].(float64)

	if ok {
		id := uint32(fid)
		ch := c.responseChans[id]
		delete(c.responseChans, id)
		ch <- obj["result"]
	} else {
		logrus.Errorf("Invalid received object during handleRecv: %T", obj["id"])
	}
}

func (c *Client) dispatcher() {
	for {
		select {
		case msg := <-c.msgs:
			c.handleSend(msg)
		case resp := <-c.resps:
			c.handleRecv(resp)
		}
	}
}

func (c *Client) read(errChan chan error) {
	decoder := json.NewDecoder(c.conn)

	for decoder.More() {
		var obj map[string]interface{}

		if err := decoder.Decode(&obj); err != nil {
			logrus.Errorf("error decoding during read: %v", err)
			errChan <- err
			continue
		}

		c.resps <- obj
	}
}
