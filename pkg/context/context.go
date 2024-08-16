/**
 * @api post myContext.
 *
 * User: yunshengzhu
 * Date: 2022/5/1
 * Time: 1:30 PM
 */
package context

import (
	c "context"
	"github.com/evanyxw/monster-go/pkg/logs"
	"github.com/evanyxw/monster-go/pkg/utils"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
)

const KeyPrefix = "metadatakeyprefix"

// metadatakeyprefixchannel
func NewContextDefault(options ...contextOptions) Context {
	myCtx := Context{}
	myCtx.Context = c.Background()
	for _, option := range options {
		option(&myCtx)
	}
	myCtx.Log = logs.Log.WithField("reqId", myCtx.reqId)
	if logs.Log.ServerName != "" {
		myCtx.Log = myCtx.Log.WithField("serverName", logs.Log.ServerName)
	}
	myCtx.kv = make(map[string]string)
	return myCtx
}

func NewContext(ctx c.Context) Context {
	myCtx := Context{}
	myCtx.Context = ctx
	myCtx.kv = make(map[string]string)
	md, _ := metadata.FromIncomingContext(ctx)
	myCtx.appId = getKey(md, "appid")
	myCtx.reqId = getKey(md, "reqid")
	myCtx.reqTime = getKey(md, "reqtime")
	myCtx.middlewareLogLevel = getKey(md, "middlewareloglevel")
	myCtx.Log = logs.Log.WithField("reqId", myCtx.reqId)
	for key, v := range md {
		if len(key) > len(KeyPrefix) {
			if utils.Substr(key, 0, len(KeyPrefix)) == KeyPrefix {
				if len(v) >= 1 {
					myCtx.kv[key] = v[0]
				}
			}
		}
	}
	if logs.Log.ServerName != "" {
		myCtx.Log = myCtx.Log.WithField("serverName", logs.Log.ServerName)
	}
	return myCtx
}

func getKey(md metadata.MD, key string) string {
	reqIds := md[key]
	reqId := ""
	if len(reqIds) > 0 {
		reqId = reqIds[0]
	}
	return reqId
}

type Context struct {
	c.Context
	appId              string
	reqId              string
	reqTime            string
	middlewareLogLevel string
	appCallBackUrl     string
	kv                 map[string]string
	Log                *log.Entry
	ChannelId          int64
}

type contextOptions func(*Context)

func WithAapId(appId string) contextOptions {
	return func(logs *Context) {
		logs.appId = appId
	}
}

func WithAppId(appId string) contextOptions {
	return func(logs *Context) {
		logs.appId = appId
	}
}

func WithReqId(reqId string) contextOptions {
	return func(logs *Context) {
		logs.reqId = reqId
	}
}

func WithReqTime(reqTime string) contextOptions {
	return func(logs *Context) {
		logs.reqTime = reqTime
	}
}

func WithMiddlewareLogLevel(middlewareLogLevel string) contextOptions {
	return func(logs *Context) {
		logs.middlewareLogLevel = middlewareLogLevel
	}
}

func WithAppCallBackUrl(appCallBackUrl string) contextOptions {
	return func(logs *Context) {
		logs.appCallBackUrl = appCallBackUrl
	}
}

func (my *Context) GetReqId() string {
	return my.reqId
}

func (my *Context) GetAppId() string {
	return my.appId
}

func (my *Context) SetAppId(appId string) {
	my.appId = appId
}

func (my *Context) GetMiddlewareLogLevel() string {
	return my.middlewareLogLevel
}

func (my *Context) GetAppCallBackUrl() string {
	return my.appCallBackUrl
}

func (my *Context) GetContext() c.Context {
	ctx := c.Background()
	md := make(map[string]string)
	md["appid"] = my.appId
	md["reqid"] = my.reqId
	md["reqtime"] = my.reqTime
	md["middlewareloglevel"] = my.middlewareLogLevel
	for k, v := range my.kv {
		md[k] = v
	}
	ctx = metadata.NewIncomingContext(ctx, metadata.New(md))
	return ctx
}

func (my *Context) SetKeyVal(key, val string) error {
	my.kv[KeyPrefix+key] = val
	return nil
}

func (my *Context) GetVal(key string) string {
	return my.kv[KeyPrefix+key]
}
