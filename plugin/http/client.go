package http

import (
	"context"
	"fmt"
	"github.com/Gitforxuyang/eva/config"
	error2 "github.com/Gitforxuyang/eva/util/error"
	"github.com/Gitforxuyang/eva/util/logger"
	"github.com/Gitforxuyang/eva/util/trace"
	"github.com/Gitforxuyang/eva/util/utils"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	METHOD_GET    HttpMethod = "GET"
	METHOD_POST   HttpMethod = "POST"
	METHOD_PUT    HttpMethod = "PUT"
	METHOD_DELETE HttpMethod = "DELETE"
	METHOD_HEAD   HttpMethod = "HEAD"
	METHOD_OPTION HttpMethod = "OPTION"
)

type HttpResp struct {
	Code    int32
	Message string
	Data    map[string]interface{}
}
type HttpMethod string

func (m HttpMethod) String() string {
	return string(m)
}

type Headers map[string]string

type DoFunc func(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) (*http.Response, error)

type EvaHttp interface {
	DoRpc(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) (map[string]interface{}, error)
	Do(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) ([]byte, error)
}

type evaHttp struct {
	cli    *http.Client
	addr   string
	conf   *config.EvaConfig
	log    logger.EvaLogger
	tracer *trace.Tracer
}

var (
	client *evaHttp
	lock   sync.Mutex
)

func GetHttpClient(name string) EvaHttp {
	lock.Lock()
	defer lock.Unlock()
	if client == nil {
		conf := config.GetConfig().GetHttp(name)
		h := new(evaHttp)
		h.addr = conf.Endpoint
		h.cli = &http.Client{
			Timeout: time.Second * time.Duration(conf.Timeout),
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   time.Second * time.Duration(conf.Timeout),
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				MaxIdleConns:          conf.MaxConn,
				MaxIdleConnsPerHost:   conf.MaxConn,
				MaxConnsPerHost:       conf.MaxConn,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
		h.conf = config.GetConfig()
		h.log = logger.GetLogger()
		h.tracer = trace.GetTracer()
		client = h
	}
	return client
}
func (m *evaHttp) DoRpc(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) (resp map[string]interface{}, err error) {
	start := time.Now()
	ctx, span, err := m.tracer.StartHttpClientSpanFromContext(ctx, fmt.Sprintf("%s_%s", uri, method))
	if err != nil {
		m.log.Error(ctx, "链路错误", logger.Fields{"err": utils.StructToMap(err)})
	}
	var statusCode int
	defer span.Finish()
	defer func() {
		span.LogFields(
			log.Object("req", utils.StructToJson(data)),
			log.Object("resp", utils.StructToJson(resp)),
			log.Object("headers", utils.StructToJson(headers)),
			//log.String("uri", uri),
			//log.String("addr", m.addr),
			//log.String("method", method.String()),
			log.String("http_type", "DoRpc"),
		)
		ext.HTTPMethod.Set(span, method.String())
		ext.HTTPUrl.Set(span, fmt.Sprintf("%s%s", m.addr, uri))
		ext.HTTPStatusCode.Set(span, uint16(statusCode))
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"))
			span.LogFields(
				log.Object("evaError", utils.StructToJson(err)),
			)
		}
	}()

	defer func() {
		if m.conf.GetLogConfig().HttpClient {
			m.log.Info(ctx, "发起的http请求", logger.Fields{
				"req":     data,
				"resp":    resp,
				"headers": headers,
				"uri":     uri,
				"addr":    m.addr,
				"method":  method,
				"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
				"err":     utils.StructToMap(err),
			})
		}
	}()
	res, err := m.do(ctx, uri, method, headers, data)
	if err != nil {
		return nil, err
	}
	statusCode = res.StatusCode
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	respData := new(HttpResp)
	err = utils.JsonToStruct(string(body), respData)
	if err != nil {
		return nil, err
	}
	//如果返回的http码不为200
	if res.StatusCode != 200 {
		return nil, error2.HttpError.SetAppId(m.conf.GetName()).SetCode(respData.Code).SetMessage(respData.Message)
	}
	if respData.Code != 0 {
		return nil, error2.HttpError.SetAppId(m.conf.GetName()).SetCode(respData.Code).SetMessage(respData.Message)
	}
	return respData.Data, nil
}

func (m *evaHttp) Do(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) (bytes []byte, err error) {
	start := time.Now()
	ctx, span, err := m.tracer.StartHttpClientSpanFromContext(ctx, fmt.Sprintf("%s_%s", uri, method))
	if err != nil {
		m.log.Error(ctx, "链路错误", logger.Fields{"err": utils.StructToMap(err)})
	}
	var statusCode int
	defer span.Finish()
	defer func() {
		span.LogFields(
			log.Object("req", utils.StructToJson(data)),
			log.Object("resp", string(bytes)),
			log.Object("headers", utils.StructToJson(headers)),
			//log.String("uri", uri),
			log.String("http_type", "Do"),
			//log.String("addr", m.addr),
		)
		ext.HTTPMethod.Set(span, method.String())
		ext.HTTPUrl.Set(span, fmt.Sprintf("%s%s", m.addr, uri))
		ext.HTTPStatusCode.Set(span, uint16(statusCode))
		if err != nil {
			ext.Error.Set(span, true)
			span.LogFields(log.String("event", "error"))
			span.LogFields(
				log.Object("evaError", utils.StructToJson(err)),
			)
		}
	}()

	defer func() {
		if m.conf.GetLogConfig().HttpClient {
			resp, _ := utils.JsonToMap(string(bytes))
			m.log.Info(ctx, "发起的http请求", logger.Fields{
				"req":     data,
				"resp":    resp,
				"headers": headers,
				"uri":     uri,
				"addr":    m.addr,
				"method":  method,
				"useTime": fmt.Sprintf("%s", time.Now().Sub(start).String()),
				"err":     utils.StructToMap(err),
			})
		}
	}()
	defer func() {
	}()
	resp, err := m.do(ctx, uri, method, headers, data)
	if err != nil {
		return nil, err
	}
	statusCode = resp.StatusCode
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	//如果返回的http码不为200
	if resp.StatusCode != 200 {
		return nil, error2.HttpError.SetAppId(m.conf.GetName()).SetCode(int32(2000 + resp.StatusCode)).SetDetail(string(body))
	}
	return body, nil
}

func (m *evaHttp) do(ctx context.Context, uri string, method HttpMethod, headers Headers, data map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", m.addr, uri)

	req, err := http.NewRequest(method.String(), url, strings.NewReader(utils.StructToJson(data)))
	if err != nil {
		return nil, err
	}
	req.Method = method.String()
	if method == METHOD_GET {
		values := req.URL.Query()
		for k, v := range data {
			vv, ok := v.(string)
			if !ok {
				return nil, error2.TypeError.SetAppId(m.conf.GetName()).SetMessage("get请求的参数值只能为string类型")
			}
			values.Add(k, vv)
		}
		req.URL.RawQuery = values.Encode()
	} else {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Add("Connection", "keep-alive")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := m.cli.Do(req)
	return resp, err
}
