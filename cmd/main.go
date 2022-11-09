package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/wtrb/webfront"
	"golang.org/x/crypto/acme/autocert"
)

var (
	httpAddr     = flag.String("http", ":http", "HTTP listen `address`")
	metricsAddr  = flag.String("metrics", "", "metrics HTTP listen `address`")
	ruleFile     = flag.String("rules", "", "rule definition `file`")
	pollInterval = flag.Duration("poll", 10*time.Second, "rule file poll `interval`")
	letsCacheDir = flag.String("letsencrypt_cache", "", "letsencrypt cache `directory` (default is to disable HTTPS)")
)

func main() {
	flag.Parse()

	s, err := webfront.New(*ruleFile, *pollInterval)
	if err != nil {
		log.Fatal(err)
	}

	if *metricsAddr != "" {
		go func() {
			log.Fatal(http.ListenAndServe(*metricsAddr, s.MetricsHandler()))
		}()
	}

	if *letsCacheDir != "" {
		m := &autocert.Manager{
			Cache:      autocert.DirCache(*letsCacheDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: s.HostPolicy,
		}
		c := tls.Config{GetCertificate: m.GetCertificate}

		l, err := tls.Listen("tcp", ":https", &c)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			log.Fatal(http.Serve(l, s))
		}()

		log.Fatal(http.ListenAndServe(*httpAddr, m.HTTPHandler(s)))
	} else {
		log.Fatal(http.ListenAndServe(*httpAddr, s))
	}
}
