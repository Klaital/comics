package filters

import (
	"context"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/klaital/comics/pkg/config"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func extractOriginatingIP(req *restful.Request) (proxyExists bool, originatingIP string) {
	ip, hasForwardedFor := req.Request.Header["X-Forwarded-For"]

	if hasForwardedFor {
		return hasForwardedFor, strings.Join(ip, " ")
	} else {
		return hasForwardedFor, strings.Split(req.Request.RemoteAddr, ":")[0]
	}
}

func GetRequestContext(req *restful.Request) context.Context {
	// Extract the context from the request attributes
	if ctx, ok := req.Attribute("ctx").(context.Context); ok {
		return ctx
	} else {
		// Lazy init the context

		// Add the RequestID
		uniqueId := uuid.NewV4()
		ctx = context.WithValue(context.Background(), "RequestID", uniqueId.String())

		req.SetAttribute("ctx", ctx)
		return ctx
	}
}
func GetContextLogger(ctx context.Context) *log.Entry {
	cfg := config.Load()
	logger := cfg.LogContext

	// Extract the RequestID from the context
	requestID, ok := ctx.Value("RequestID").(string)
	if ok && len(requestID) > 0 {
		logger = cfg.LogContext.WithField("RequestID", requestID)
	}

	return logger
}


func GetRequestLogger(req *restful.Request) *log.Entry {
	// Extract the context from the request attributes
	ctx := GetRequestContext(req)
	return GetContextLogger(ctx)
}


func RequestLogFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	// Get the decorated logger for the request
	logger := GetRequestLogger(req)

	// W3C Trace Context Traceparent header
	// The traceparent header represents the incoming request in a tracing system in a common format, understood by all vendors.
	traceparent := req.HeaderParameter("traceparent")

	// Extract the originating IP from the header provided by the load balancer
	_, originatingIP := extractOriginatingIP(req)
	fields := log.Fields{
		"operation":     "RequestLogger",
		"traceparent":   traceparent,
		"originatingIP": originatingIP,
		"time":          time.Now().Format("02/Jan/2006:15:04:05.000 -0700"),
		"method":        req.Request.Method,
		"requestURI":    req.Request.URL.RequestURI(),
		"proto":         req.Request.Proto,
		"statusCode":    resp.StatusCode(),
		"contentLength": resp.ContentLength(),
		"user-agent":    req.HeaderParameter("user-agent"),
		// TODO: log the request & response body as well if in Debug mode
	}
	logMsgNCSACLF := fmt.Sprintf("%s - - [%s] \"%s %s %s\" %d %d",
		originatingIP,
		time.Now().Format("02/Jan/2006:15:04:05.000 -0700"),
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
	)

	// TODO: once the service is working, hide this for healthcheck requests
	logger.WithFields(fields).Info(logMsgNCSACLF)
	chain.ProcessFilter(req, resp)
}
