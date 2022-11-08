package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/wtrb/webfront"
)

var (
	httpAddr     = flag.String("http", ":http", "HTTP listen `address`")
	ruleFile     = flag.String("rules", "", "rule definition `file`")
	pollInterval = flag.Duration("poll", 10*time.Second, "rule file poll `interval`")
)

func main() {
	flag.Parse()

	s, err := webfront.New(*ruleFile, *pollInterval)
	if err != nil {
		log.Fatal(err)
	}

	if err := http.ListenAndServe(*httpAddr, s); err != nil {
		log.Fatal(err)
	}
}
