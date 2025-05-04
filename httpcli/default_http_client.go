package httpcli

import (
	"context"
	"errors"
	"fmt"
	"github.com/shadowors/goex/v2/logger"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	Cli     IHttpClient
	gHeader map[string]string
)

func init() {
	gHeader = make(map[string]string, 2)
	Cli = NewDefaultHttpClient()
}

type DefaultHttpClient struct {
	cli     *http.Client
	timeout time.Duration
}

func NewDefaultHttpClient() *DefaultHttpClient {
	cli := new(DefaultHttpClient)
	cli.init()
	return cli
}

func (cli *DefaultHttpClient) init() {
	logger.Info("[http utils] setup default http client")

	cli.timeout = 5 * time.Second
	cli.cli = &http.Client{
		Timeout: cli.timeout,
		Transport: &http.Transport{
			IdleConnTimeout:       time.Minute,
			TLSHandshakeTimeout:   cli.timeout,
			ResponseHeaderTimeout: cli.timeout,
			MaxConnsPerHost:       64,
			MaxIdleConnsPerHost:   64,
			MaxIdleConns:          128,
			ForceAttemptHTTP2:     true,
		},
	}
}

func (cli *DefaultHttpClient) SetHeaders(key, value string) {
	gHeader[key] = value
}

func (cli *DefaultHttpClient) SetTimeout(sec int64) {
	timeout := time.Duration(sec) * time.Second
	cli.timeout = timeout
	logger.Infof("[DefaultHttpClient] http(s) timeout: %s", timeout.String())
	cli.cli.Timeout = timeout
	trans := cli.cli.Transport.(*http.Transport)
	trans.ResponseHeaderTimeout = timeout
	trans.TLSHandshakeTimeout = timeout
	trans.ExpectContinueTimeout = timeout
}

func (cli *DefaultHttpClient) SetProxy(proxy string) error {
	proxyUrl, err := url.Parse(proxy)
	if err != nil {
		logger.Warnf("[DefaultHttpClient] parse proxy url err: %s", err.Error())
		return err
	}
	logger.Infof("[DefaultHttpClient] http(s) proxy url: %s", proxy)
	trans := cli.cli.Transport.(*http.Transport)
	trans.Proxy = func(request *http.Request) (*url.URL, error) {
		return proxyUrl, nil
	}
	return nil
}

func (cli *DefaultHttpClient) DoRequest(method, rqUrl string, reqBody string, headers map[string]string) (data []byte, err error) {
	logger.Debugf("[DefaultHttpClient] [%s] request url: %s", method, rqUrl)

	reqTimeoutCtx, cancelFn := context.WithTimeout(context.TODO(), cli.timeout)
	defer cancelFn()

	req, err := http.NewRequestWithContext(reqTimeoutCtx, method, rqUrl, strings.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %w", err)
	}

	//append global http header
	for k, v := range gHeader {
		req.Header.Set(k, v)
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	resp, err := cli.cli.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error("[DefaultHttpClient] close response body error:", err.Error())
		}
	}(resp.Body)

	bodyData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}

	if resp.StatusCode != 200 {
		return bodyData, errors.New(resp.Status)
	}

	return bodyData, nil
}
