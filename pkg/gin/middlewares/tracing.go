package middlewares

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const defaultComponentName = "net/http"

type trOptions struct {
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

func (o trOptions) GetComponentName() string {
	componentName := o.componentName
	if componentName == "" {
		return defaultComponentName
	}
	return componentName
}

// TracingOption controls the behavior of the Middleware.
type TracingOption func(*trOptions)

// OperationNameFunc returns a TracingOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *http.Request) string) TracingOption {
	return func(options *trOptions) {
		options.opNameFunc = f
	}
}

// MWComponentName returns a TracingOption that sets the component name
// for the server-side span.
func TracingComponentName(componentName string) TracingOption {
	return func(options *trOptions) {
		options.componentName = componentName
	}
}

// MWSpanObserver returns a TracingOption that observe the span
// for the server-side span.
func TracingSpanObserver(f func(span opentracing.Span, r *http.Request)) TracingOption {
	return func(options *trOptions) {
		options.spanObserver = f
	}
}

// MWURLTagFunc returns a TracingOption that uses given function f
// to set the span's http.url tag. Can be used to change the default
// http.url tag, eg to redact sensitive information.
func TracingURLTagFunc(f func(u *url.URL) string) TracingOption {
	return func(options *trOptions) {
		options.urlTagFunc = f
	}
}

// Middleware is a gin native version of the equivalent middleware in:
//
//	https://github.com/opentracing-contrib/go-stdlib/
func Tracing(options ...TracingOption) gin.HandlerFunc {
	opts := trOptions{
		opNameFunc: func(r *http.Request) string {
			return r.Method + " " + r.URL.Path
		},
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}

	return func(c *gin.Context) {
		tr := opentracing.GlobalTracer()
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)
		componentName := opts.GetComponentName()
		op := componentName + " " + opts.opNameFunc(c.Request)
		sp := tr.StartSpan(op, ext.RPCServerOption(ctx))
		ext.HTTPMethod.Set(sp, c.Request.Method)
		ext.HTTPUrl.Set(sp, opts.urlTagFunc(c.Request.URL))
		opts.spanObserver(sp, c.Request)

		ext.Component.Set(sp, componentName)
		c.Request = c.Request.WithContext(
			opentracing.ContextWithSpan(c.Request.Context(), sp))

		c.Next()

		ext.HTTPStatusCode.Set(sp, uint16(c.Writer.Status()))
		sp.Finish()
	}
}
