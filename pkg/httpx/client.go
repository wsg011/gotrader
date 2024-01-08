package httpx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	maxIdleConnPerHost = 200
	idleConnTimeout    = 90
	requestTimeout     = 4
	dialTimeout        = 5
	keepAlive          = 30
)

var defaultConfig = Config{
	// 每个host的idle状态的最大连接数目
	MaxIdleConnsPerHost: maxIdleConnPerHost,
	// 连接保持idle状态的最大时间，超时关闭conn
	IdleConnTimeout: idleConnTimeout * time.Second,
	// 请求超时时间，指连接时间，任何重定向时间和读取响应时间总和
	RequestTimeout: requestTimeout * time.Second,
	// 拨号等待连接完成的最大时间
	DialTimeout: dialTimeout * time.Second,
	// 心跳时间
	KeepAlive: keepAlive * time.Second,
}

type Client struct {
	*http.Client
}

type Config struct {
	// 所有host的idle状态的最大连接数目
	MaxIdleConns int
	// 每个host的idle状态的最大连接数目
	MaxIdleConnsPerHost int
	// 每个host上的最大连接数目，含dialing/active/idle状态的conn
	MaxConnsPerHost int
	// 连接保持idle状态的最大时间，超时关闭conn
	IdleConnTimeout time.Duration
	// 请求超时时间，指连接时间，任何重定向时间和读取响应时间总和
	RequestTimeout time.Duration
	// 拨号等待连接完成的最大时间
	DialTimeout time.Duration
	// 心跳时间
	KeepAlive time.Duration
}

type Request struct {
	Url    string
	Head   map[string]string
	Method string
	Body   []byte
}

func GetDefaultConfig() Config {
	return defaultConfig
}

func NewClient() Client {
	return NewClientWithConfig(defaultConfig)
}

func NewClientWithConfig(config Config) Client {
	c := &http.Client{
		Timeout: config.RequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   config.DialTimeout,
				KeepAlive: config.KeepAlive,
			}).DialContext,
			MaxIdleConns:        config.MaxIdleConns,
			MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
			MaxConnsPerHost:     config.MaxConnsPerHost,
			IdleConnTimeout:     config.IdleConnTimeout,
		},
	}
	return Client{c}
}

func (c *Client) Get(url string) (*[]byte, *http.Response, error) {
	res, err := c.Client.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("请求%s出错:%v", url, err)
	}
	defer res.Body.Close()
	body, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return nil, nil, fmt.Errorf("读取%v错误:%v", res, err2)
	}
	return &body, res, nil
}

func (c *Client) Post(url string, body interface{}) (*[]byte, *http.Response, error) {
	byteBody, _ := json.Marshal(body)
	res, err := c.Client.Post(url, "application/json", bytes.NewBuffer(byteBody))
	if err != nil {
		return nil, nil, fmt.Errorf("请求%s出错:%v", url, err)
	}
	defer res.Body.Close()
	resBody, err2 := io.ReadAll(res.Body)
	if err2 != nil {
		return nil, nil, fmt.Errorf("读取%v错误:%v", res, err2)
	}
	return &resBody, res, nil
}

func (c *Client) Request(args *Request) (*[]byte, *http.Response, error) {
	url, head, method, body := args.Url, args.Head, args.Method, args.Body
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, fmt.Errorf("创建http请求出错:%v", err)
	}
	if head == nil {
		req.Header.Set("Content-Type", "application/json")
	} else {
		for key, val := range head {
			req.Header.Set(key, val)
		}
	}
	resp, err2 := c.Do(req)
	if err2 != nil {
		return nil, nil, fmt.Errorf("%s请求%s出错:%v", method, url, err2)
	}
	defer resp.Body.Close()
	rspBody, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		return nil, nil, fmt.Errorf("读取%v错误:%v", resp, err3)
	}
	return &rspBody, resp, nil
}
