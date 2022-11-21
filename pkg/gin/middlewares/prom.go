package middlewares

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultPath = "/metrics"

type pmapb struct {
	sync.RWMutex
	values map[string]bool
}

// Prometheus contains the metrics gathered by the instance and its path
type Prometheus struct {
	reqCnt               *prometheus.CounterVec
	reqDur, reqSz, resSz *prometheus.HistogramVec

	MetricsPath string
	Ignored     pmapb
}

// Ignore is used to disable instrumentation on some routes
func (p *Prometheus) Ignore(paths ...string) {
	p.Ignored.Lock()
	defer p.Ignored.Unlock()
	for _, path := range paths {
		p.Ignored.values[path] = true
	}
}

// NewPromMiddleware will initialize a new Prometheus instance with the given options.
// If no options are passed, sane defaults are used.
// If a router is passed using the Engine() option, this instance will
// automatically bind to it.
func NewPromMiddleware() *Prometheus {
	p := &Prometheus{
		MetricsPath: defaultPath,
	}
	p.Ignored.values = make(map[string]bool)

	p.register()

	return p
}

func (p *Prometheus) register() {
	p.reqCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"endpoint", "code", "method", "path"},
	)
	prometheus.MustRegister(p.reqCnt)

	p.reqDur = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "The HTTP request latencies in seconds.",
		},
		[]string{"endpoint", "code", "method", "path"},
	)
	prometheus.MustRegister(p.reqDur)

	p.reqSz = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "The HTTP request sizes in bytes.",
			Buckets: []float64{1, 5, 10, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000},
		},
		[]string{"endpoint", "method", "path"},
	)
	prometheus.MustRegister(p.reqSz)

	p.resSz = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "The HTTP response sizes in bytes.",
			Buckets: []float64{1, 5, 10, 50, 100, 250, 500, 1000, 2500, 5000, 10000, 25000, 50000},
		},
		[]string{"endpoint", "method", "path"},
	)
	prometheus.MustRegister(p.resSz)
}

// Instrument is a gin middleware that can be used to generate metrics for a
// single handler
func (p *Prometheus) Instrument(endpoint string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		reqSz := computeApproximateRequestSize(c.Request)

		c.Next()

		path := c.FullPath()
		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		p.reqCnt.WithLabelValues(endpoint, status, c.Request.Method, path).Inc()
		p.reqDur.WithLabelValues(endpoint, status, c.Request.Method, path).Observe(elapsed)
		p.reqSz.WithLabelValues(endpoint, c.Request.Method, path).Observe(float64(reqSz))
		p.resSz.WithLabelValues(endpoint, c.Request.Method, path).Observe(resSz)
	}
}

// Use is a method that should be used if the engine is set after middleware
// initialization
func (p *Prometheus) Use(e *gin.Engine) {
	e.GET(p.MetricsPath, prometheusHandler())
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}
