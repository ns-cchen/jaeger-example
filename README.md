1. Initialize jaeger docker
```shell
docker run --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 14250:14250 -p 9411:9411 jaegertracing/all-in-one
```

2. Set up environment variables
```shell
JAEGER_AGENT_HOST=127.0.0.1
JAEGER_ENDPOINT=http://127.0.0.1:14268/api/traces
JAEGER_SERVICE_NAME=microservice
```

3. Launch server
```shell
 go run main.go
```

4. Send an API
```shell
curl http://localhost:5678/server1  
```

5. Go to http://localhost:16686/, you can see the call path
