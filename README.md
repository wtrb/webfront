Webfront is an HTTP server and reverse proxy.

https://go.dev/talks/2012/simple.slide#29
https://github.com/nf/webfront

Prothemeus: https://pkg.go.dev/github.com/prometheus/client_golang/prometheus

Let's Encrypt

Run:
```
go run main.go -http :80 -rules ./rules.json -poll 10s -metrics :81
```
