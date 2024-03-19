package server

import (
	"net/http"

	"github.com/netskope/go-kestrel/pkg/server"
	"jaeger-example/client"
)

func Run() {
	go func() {
		server1, _ := server.New(server.WithListenAddr("localhost:5678"), server.TracingMiddleware())
		server1.ChiRouter().Get("/server1", func(writer http.ResponseWriter, request *http.Request) {
			body := client.Get(request, "client", "server2", 9101)
			writer.WriteHeader(200)
			_, _ = writer.Write([]byte(body))
		})
		_ = server1.Run()
	}()

	go func() {
		server2, _ := server.New(server.WithListenAddr("localhost:9101"), server.TracingMiddleware())
		server2.ChiRouter().Get("/server2", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(200)
			_, _ = writer.Write([]byte("hello"))
		})
		_ = server2.Run()
	}()

	select {}
}

//
//func getGreeting(tracer opentracing.Tracer, parentSpan opentracing.Span, req *http.Request) {
//	childSpan := tracer.StartSpan(
//		"Call API",
//		opentracing.ChildOf(parentSpan.Context()),
//	)
//	ext.SpanKindRPCClient.Set(childSpan)
//	ext.HTTPMethod.Set(childSpan, "GET")
//	ext.HTTPUrl.Set(childSpan, req.URL.String())
//
//	_ = tracer.Inject(childSpan.Context(),
//		opentracing.HTTPHeaders,
//		opentracing.HTTPHeadersCarrier(req.Header))
//	println(getBody(req))
//
//	defer childSpan.Finish()
//}
//
//func getRequest(path string, port int) *http.Request {
//
//	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%v/%s", port, path), nil)
//	if err != nil {
//		fmt.Println("Error creating request:", err)
//		return nil
//	}
//
//	return req
//}
//
//func getBody(req *http.Request) string {
//	client := &http.Client{}
//
//	resp, err := client.Do(req)
//	if err != nil {
//		fmt.Println("Error sending request:", err)
//		return ""
//	}
//	defer func(Body io.ReadCloser) {
//		_ = Body.Close()
//	}(resp.Body)
//
//	body, err := io.ReadAll(resp.Body)
//	if err != nil {
//		fmt.Println("Error reading response body:", err)
//		return ""
//	}
//
//	return string(body)
//}
//
//func createTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
//	var cfg = jaegercfg.Configuration{
//		ServiceName: serviceName,
//		Sampler: &jaegercfg.SamplerConfig{
//			Type:  jaeger.SamplerTypeConst,
//			Param: 1,
//		},
//		Reporter: &jaegercfg.ReporterConfig{
//			LogSpans:          true,
//			CollectorEndpoint: "http://localhost:14268/api/traces",
//		},
//	}
//
//	jLogger := jaegerlog.StdLogger
//	tracer, closer, err := cfg.NewTracer(
//		jaegercfg.Logger(jLogger),
//	)
//	return tracer, closer, err
//}
