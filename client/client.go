package client

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"jaeger-example/tracing"
)

func Get(r *http.Request, serviceName string, path string, port int) string {
	tracer, closer, _ := createTracer(serviceName)
	defer func(closer io.Closer) {
		_ = closer.Close()
	}(closer)

	var span opentracing.Span
	if r == nil {
		span = tracer.StartSpan("Start")
	} else {
		span = tracing.StartSpanFromRequest(tracer, r, "Call API")
	}
	defer span.Finish()

	var req *http.Request
	if r == nil {
		req = getRequest(path, port, nil)
	} else {
		req = getRequest(path, port, r.Header)
	}

	ext.SpanKindRPCClient.Set(span)
	ext.HTTPMethod.Set(span, "GET")
	ext.HTTPUrl.Set(span, req.URL.String())

	_ = tracing.Inject(span, req)
	return getBody(req)
}

func getRequest(path string, port int, header http.Header) *http.Request {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/%s", port, path), nil)
	if header != nil {
		for key, values := range header {
			req.Header.Add(key, values[0])
		}
	}

	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	return req
}

func getBody(req *http.Request) string {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return ""
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

func createTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}

	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	return tracer, closer, err
}
