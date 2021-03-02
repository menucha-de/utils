package utils

import (
	"net/rpc"
	"reflect"
	"sync"
)

type Client struct {
	Mutex        sync.Mutex
	rpcClient    *rpc.Client
	ServerAdress string
}

func (c *Client) Call(serviceMethod string, args interface{}, reply interface{}) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	var err error
	rpcClient := c.rpcClient
	c.Mutex.Unlock()
	if rpcClient != nil {
		err = rpcClient.Call(serviceMethod, args, reply)
		c.Mutex.Lock()
		if err != nil {
			//Check if call is caused by disconnect/any other server error
			if err == rpc.ErrShutdown || reflect.TypeOf(err) == reflect.TypeOf((*rpc.ServerError)(nil)).Elem() {
				rpcClient.Close()
				c.rpcClient, rpcClient = nil, nil
				c.Mutex.Unlock()
			} else {
				return err
			}
		}
	}

	//Re-/Initialize connection for call
	if rpcClient == nil {
		c.Mutex.Lock()
		if c.rpcClient == nil {
			c.rpcClient, err = rpc.DialHTTP("tcp", c.ServerAdress)
			if err != nil {
				return err
			}
		}
		rpcClient = c.rpcClient
		c.Mutex.Unlock()
		err = rpcClient.Call(serviceMethod, args, reply)
		c.Mutex.Lock()
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.rpcClient.Close()
}
