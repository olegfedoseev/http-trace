package trace

import (
	"context"
	"fmt"
	"net/http/httptrace"
	"time"
)

// Result will hold results of request tracing
type Result struct {
	URL              string
	DNSLookup        time.Duration
	TCPConnection    time.Duration
	ServerProcessing time.Duration
	ContentTransfer  time.Duration
	Total            time.Duration
}

func (r Result) String() string {
	return fmt.Sprintf("%s [DNS: %v TCP: %v Server: %v Transfer: %v Total: %v]",
		r.URL, r.DNSLookup, r.TCPConnection, r.ServerProcessing, r.ContentTransfer, r.Total)
}

// Trace will hold time of trace events of request lifecycle
type Trace struct {
	URL                      string
	Created                  time.Time
	RequestStart             time.Time
	DNSStartTime             time.Time
	DNSDoneTime              time.Time
	GotConnTime              time.Time
	GotFirstResponseByteTime time.Time
	BodyReadTime             time.Time
}

// NewTrace will create new Trace and context with httptrace.ClientTrace
func NewTrace(url string) (*Trace, context.Context) {
	clientTrace := &Trace{
		URL:     url,
		Created: time.Now(),
	}
	trace := &httptrace.ClientTrace{
		DNSStart:             clientTrace.DNSStart,
		DNSDone:              clientTrace.DNSDone,
		ConnectStart:         clientTrace.ConnectStart,
		GotConn:              clientTrace.GotConn,
		GotFirstResponseByte: clientTrace.GotFirstResponseByte,
	}

	return clientTrace, httptrace.WithClientTrace(context.Background(), trace)
}

// GetResult will return Result struct with easy to use fields with measurements
func (t *Trace) GetResult() Result {
	if t.RequestStart.IsZero() {
		t.RequestStart = t.Created
	}
	return Result{
		URL:              t.URL,
		DNSLookup:        t.DNSDoneTime.Sub(t.DNSStartTime),
		TCPConnection:    t.GotConnTime.Sub(t.RequestStart),
		ServerProcessing: t.GotFirstResponseByteTime.Sub(t.GotConnTime),
		ContentTransfer:  t.BodyReadTime.Sub(t.GotFirstResponseByteTime),
		Total:            t.BodyReadTime.Sub(t.RequestStart),
	}
}

// DNSStart will be used in httptrace.ClientTrace DNSStart
func (t *Trace) DNSStart(_ httptrace.DNSStartInfo) {
	t.RequestStart = time.Now()
	t.DNSStartTime = time.Now()
}

// DNSDone will be used in httptrace.ClientTrace DNSDone
func (t *Trace) DNSDone(_ httptrace.DNSDoneInfo) {
	t.DNSDoneTime = time.Now()
}

// ConnectStart will be used in httptrace.ClientTrace ConnectStart
func (t *Trace) ConnectStart(_, _ string) {
	if t.RequestStart.IsZero() {
		t.RequestStart = time.Now()
	}
}

// GotFirstResponseByte will be used in httptrace.ClientTrace GotFirstResponseByte
func (t *Trace) GotFirstResponseByte() {
	t.GotFirstResponseByteTime = time.Now()
}

// GotConn will be used in httptrace.ClientTrace GotConn
func (t *Trace) GotConn(info httptrace.GotConnInfo) {
	t.GotConnTime = time.Now()
}
