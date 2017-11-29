package trace

import (
	"io"
	"net/http"
	"time"
)

// spyBodyReader is helper to detect when response body is read
type spyBodyReader struct {
	io.ReadCloser
	// Pointer to var in trace struct
	bodyReadTime *time.Time
}

func (br spyBodyReader) Read(p []byte) (n int, err error) {
	n, err = br.ReadCloser.Read(p)
	if err != nil {
		// If there was some error, including io.EOF than we consider it as the end
		*br.bodyReadTime = time.Now()
	}
	return
}

// DoRequestWithTrace is helper to make a request with trace
func DoRequestWithTrace(client *http.Client, req *http.Request) (*Trace, *http.Response, error) {
	clientTrace, ctx := NewTrace(req.URL.String())
	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if resp != nil {
		resp.Body = spyBodyReader{
			ReadCloser:   resp.Body,
			bodyReadTime: &clientTrace.BodyReadTime,
		}
	}

	return clientTrace, resp, err
}
