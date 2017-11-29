# http-trace
Helper for golang's httptrace.ClientTrace

# How to use

```
import (
    trace "github.com/olegfedoseev/http-trace"
)

// ...

    req, _ := http.NewRequest(...)
    httpClient := &http.Client{}

    // insted of 
    // resp, err := httpClient.Do(req)
    // you need to do
    clientTrace, resp, err := trace.DoRequestWithTrace(httpClient, req)
    if err != nil {
        // handle err...
    }
    defer func() {
        resp.Body.Close()
        fmt.Printf("[TRACE] %s\n", clientTrace.GetResult())
    }()
    
// ...
```
And you'll get something like than:

```
[TRACE] http://example.com [DNS: 0s TCP: 7.778836ms Server: 85.851965ms Transfer: 136.281µs Total: 93.767082ms]
```
Or you can format it the way you like, `clientTrace.GetResult()` gives you access to all data.
