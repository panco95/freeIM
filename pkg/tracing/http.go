package tracing

import (
	"crypto/tls"
	"net/http/httptrace"
	"net/textproto"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

func NewClientTrace(span opentracing.Span) *httptrace.ClientTrace {
	trace := &clientTrace{span: span}
	return &httptrace.ClientTrace{
		GetConn:              trace.GetConn,
		GotConn:              trace.GotConn,
		PutIdleConn:          trace.PutIdleConn,
		GotFirstResponseByte: trace.GotFirstResponseByte,
		Got100Continue:       trace.Got100Continue,
		Got1xxResponse:       trace.Got1xxResponse,
		DNSStart:             trace.DNSStart,
		DNSDone:              trace.DNSDone,
		ConnectStart:         trace.ConnectStart,
		ConnectDone:          trace.ConnectDone,
		TLSHandshakeStart:    trace.TLSHandshakeStart,
		TLSHandshakeDone:     trace.TLSHandshakeDone,
		WroteRequest:         trace.WroteRequest,
		WroteHeaders:         trace.WroteHeaders,
		Wait100Continue:      trace.Wait100Continue,
	}
}

// clientTrace holds a reference to the Span and
// provides methods used as ClientTrace callbacks
type clientTrace struct {
	span opentracing.Span
}

func (h *clientTrace) GetConn(hostPort string) {
	h.span.LogFields(
		log.String("event", "GetConn"),
		log.String("hostPort", hostPort),
	)
}

func (h *clientTrace) GotConn(info httptrace.GotConnInfo) {
	h.span.SetTag("net/http.reused", info.Reused)
	h.span.SetTag("net/http.was_idle", info.WasIdle)
	h.span.LogFields(log.String("event", "GotConn"))
}

func (h *clientTrace) PutIdleConn(err error) {
	h.span.LogFields(log.String("event", "PutIdleConn"))
}

func (h *clientTrace) GotFirstResponseByte() {
	h.span.LogFields(log.String("event", "GotFirstResponseByte"))
}

func (h *clientTrace) Got100Continue() {
	h.span.LogFields(log.String("event", "Got100Continue"))
}

func (h *clientTrace) Got1xxResponse(code int, header textproto.MIMEHeader) error {
	h.span.LogFields(log.String("event", "Got1xxResponse"))
	return nil
}

func (h *clientTrace) DNSStart(info httptrace.DNSStartInfo) {
	h.span.LogFields(
		log.String("event", "DNSStart"),
		log.Object("host", info.Host),
	)
}

func (h *clientTrace) DNSDone(info httptrace.DNSDoneInfo) {
	fields := []log.Field{log.String("event", "DNSDone")}
	for _, addr := range info.Addrs {
		fields = append(fields, log.String("addr", addr.String()))
	}
	if info.Err != nil {
		fields = append(fields, log.Error(info.Err))
	}
	h.span.LogFields(fields...)
}

func (h *clientTrace) ConnectStart(network, addr string) {
	h.span.LogFields(
		log.String("event", "ConnectStart"),
		log.String("network", network),
		log.String("addr", addr),
	)
}

func (h *clientTrace) ConnectDone(network, addr string, err error) {
	if err != nil {
		h.span.LogFields(
			log.String("message", "ConnectDone"),
			log.String("network", network),
			log.String("addr", addr),
			log.String("event", "error"),
			log.Error(err),
		)
	} else {
		h.span.LogFields(
			log.String("event", "ConnectDone"),
			log.String("network", network),
			log.String("addr", addr),
		)
	}
}

func (h *clientTrace) TLSHandshakeStart() {
	h.span.LogFields(log.String("event", "TLSHandshakeStart"))
}

func (h *clientTrace) TLSHandshakeDone(tls.ConnectionState, error) {
	h.span.LogFields(log.String("event", "TLSHandshakeDone"))
}

func (h *clientTrace) WroteHeaderField(key string, value []string) {
	h.span.LogFields(log.String("event", "WroteHeaderField"))
}

func (h *clientTrace) WroteHeaders() {
	h.span.LogFields(log.String("event", "WroteHeaders"))
}

func (h *clientTrace) Wait100Continue() {
	h.span.LogFields(log.String("event", "Wait100Continue"))
}

func (h *clientTrace) WroteRequest(info httptrace.WroteRequestInfo) {
	if info.Err != nil {
		h.span.LogFields(
			log.String("message", "WroteRequest"),
			log.String("event", "error"),
			log.Error(info.Err),
		)
		ext.Error.Set(h.span, true)
	} else {
		h.span.LogFields(log.String("event", "WroteRequest"))
	}
}
