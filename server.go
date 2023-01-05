/*
webfront is an HTTP server and reverse proxy.
It reads a JSON-formatted rule file like this:

	[
		{"Host": "example.com", "Serve": "/var/www"},
		{"Host": "example.org", "Forward": "localhost:8080"}
	]

For all requests to the host example.com (or any name ending in
".example.com") it serves files from the /var/www directory.
For requests to example.org, it forwards the request to the HTTP
server listening on localhost port 8080.
Usage of webfront:

	-http address
	  	HTTP listen address (default ":http")
	-letsencrypt_cache directory
	  	letsencrypt cache directory (default is to disable HTTPS)
	-poll interval
	  	rule file poll interval (default 10s)
	-rules file
	  	rule definition file

webfront was written by Andrew Gerrand <adg@golang.org>
*/
package webfront

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// New constructs a Server that reads rules from file with a period
// specified by poll.
func New(file string, pool time.Duration) (*Server, error) {
	s := new(Server)
	if err := s.loadRules(file); err != nil {
		return nil, err
	}

	go s.refreshRules(file, pool)

	return s, nil
}

// Server implements an http.Handler that acts as either a reverse proxy or
// a simple file server, as determined by a rule set.
type Server struct {
	mu      sync.RWMutex // guards the fields below
	lastMod time.Time    // rules' last modified time
	rules   []*Rule
}

// ServeHTTP matches the Request with a Rule and, if found, serves the
// request with the Rule's handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h := s.handler(r); h != nil {
		h.ServeHTTP(w, r)
		return
	}

	http.Error(w, "Not found.", http.StatusNotFound)
}

// handler returns the appropriate Handler for the given Request,
// or nil if none found.
func (s *Server) handler(req *http.Request) http.Handler {
	s.mu.RLock()
	defer s.mu.RUnlock()

	h := req.Host
	// Some clients include a port in the request host; strip it.
	if i := strings.Index(h, ":"); i >= 0 {
		h = h[:i]
	}
	for _, r := range s.rules {
		if h == r.Host || strings.HasSuffix(h, "."+r.Host) {
			hitCounter.With(prometheus.Labels{"host": r.Host}).Inc()

			voltageGauge.
				With(
					prometheus.Labels{
						"country": "sg",
						"site":    sites[rand.Intn(len(sites))],
						"pv":      pvs[rand.Intn(len(pvs))],
						"iv":      inverters[rand.Intn(len(inverters))],
					},
				).
				Set(float64(rand.Intn(10)))
				// Set(rand.Float64() * 100)

			return r.handler
		}
	}
	return nil
}

var (
	sites     = [...]string{"rangoon", "rangaan"}
	pvs       = [...]string{"pv A", "pv B"}
	inverters = [...]string{"iv1", "iv2", "iv3"}
)

// refreshRules polls file periodically and refreshes the Server's rule
// set if the file has been modified.
func (s *Server) refreshRules(file string, poll time.Duration) {
	for range time.Tick(poll) {
		if err := s.loadRules(file); err != nil {
			log.Println(err)
		}
	}
}

// loadRules tests whether file has been modified since its last invocation
// and, if so, loads the rule set from file.
func (s *Server) loadRules(file string) error {
	fi, err := os.Stat(file)
	if err != nil {
		return err
	}
	mtime := fi.ModTime()
	if !mtime.After(s.lastMod) && s.rules != nil {
		return nil // no change
	}

	rules, err := parseRules(file)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastMod = mtime
	s.rules = rules

	return nil
}

// hostPolicy implements autocert.HostPolicy by consulting
// the rules list for a matching host name.
func (s *Server) HostPolicy(ctx context.Context, host string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.rules {
		if host == r.Host || host == "www."+r.Host {
			return nil
		}
	}

	return fmt.Errorf("unrecognized host %q", host)
}
