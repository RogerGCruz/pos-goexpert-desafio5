package steps

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	src "github.com/RogerGCruz/pos-goexpert-desafio5/src"

	"github.com/cucumber/godog"
)

type loadTestFeature struct {
	server *httptest.Server
	output string
}

func (f *loadTestFeature) aTestHTTPServerThatReturns(status int) error {
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte("ok"))
	}))
	return nil
}

func (f *loadTestFeature) aTestHTTPServerThatCycles(spec string) error {
	parts := strings.Split(spec, ",")
	statuses := make([]int, 0, len(parts))
	for _, p := range parts {
		s := strings.TrimSpace(p)
		v, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("invalid status in cycle: %s", s)
		}
		statuses = append(statuses, v)
	}
	var cnt int64
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := atomic.AddInt64(&cnt, 1) - 1
		code := statuses[int(i)%len(statuses)]
		w.WriteHeader(code)
		_, _ = w.Write([]byte("ok"))
	}))
	return nil
}

func (f *loadTestFeature) aTestHTTPServerThatSleepsSeconds(seconds int) error {
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(seconds) * time.Second)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	return nil
}

func (f *loadTestFeature) aTestHTTPServerThatSleepsMsThenResponds(ms int, status int) error {
	f.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(ms) * time.Millisecond)
		w.WriteHeader(status)
		_, _ = w.Write([]byte("ok"))
	}))
	return nil
}

func (f *loadTestFeature) iRunTheCLIWithServerURLRequestsAndConcurrency(requests, concurrency int) error {
	var buf bytes.Buffer
	args := []string{
		fmt.Sprintf("--url=%s", f.server.URL),
		fmt.Sprintf("--requests=%d", requests),
		fmt.Sprintf("--concurrency=%d", concurrency),
	}
	code := src.Run(args, &buf)
	f.output = buf.String()
	if f.server != nil {
		f.server.Close()
	}
	if code != 0 {
		return fmt.Errorf("loadtest exited with code %d, output: %s", code, f.output)
	}
	return nil
}

func (f *loadTestFeature) theReportShowsTotalRequests(expected int) error {
	re := regexp.MustCompile(`Total de requests: (\d+)`)
	m := re.FindStringSubmatch(f.output)
	if len(m) < 2 {
		return fmt.Errorf("total requests not found in output")
	}
	v, _ := strconv.Atoi(m[1])
	if v != expected {
		return fmt.Errorf("expected total %d got %d", expected, v)
	}
	return nil
}

func (f *loadTestFeature) theReportShows200OK(expected int) error {
	re := regexp.MustCompile(`200 OK: (\d+)`)
	m := re.FindStringSubmatch(f.output)
	if len(m) < 2 {
		return fmt.Errorf("200 OK count not found in output")
	}
	v, _ := strconv.Atoi(m[1])
	if v != expected {
		return fmt.Errorf("expected 200 OK %d got %d", expected, v)
	}
	return nil
}

func (f *loadTestFeature) theReportShowsStatusCount(code, expected int) error {
	re := regexp.MustCompile(fmt.Sprintf("%d: (\\d+)", code))
	m := re.FindStringSubmatch(f.output)
	if len(m) < 2 {
		return fmt.Errorf("status %d count not found in output", code)
	}
	v, _ := strconv.Atoi(m[1])
	if v != expected {
		return fmt.Errorf("expected status %d count %d got %d", code, expected, v)
	}
	return nil
}

func (f *loadTestFeature) theReportTotalAtLeastMs(minMs int) error {
	re := regexp.MustCompile(`Tempo total: ([^\n]+)`)
	m := re.FindStringSubmatch(f.output)
	if len(m) < 2 {
		return fmt.Errorf("tempo total not found in output")
	}
	d, err := time.ParseDuration(m[1])
	if err != nil {
		return fmt.Errorf("cannot parse duration: %v", err)
	}
	if d < time.Duration(minMs)*time.Millisecond {
		return fmt.Errorf("expected total >= %dms got %s", minMs, d)
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	f := &loadTestFeature{}
	ctx.Step(`^a test HTTP server that returns (\d+)$`, func(code int) error { return f.aTestHTTPServerThatReturns(code) })
	ctx.Step(`^a test HTTP server that cycles through (.+)$`, func(spec string) error { return f.aTestHTTPServerThatCycles(spec) })
	ctx.Step(`^a test HTTP server that sleeps (\d+)s before responding$`, func(s int) error { return f.aTestHTTPServerThatSleepsSeconds(s) })
	ctx.Step(`^a test HTTP server that sleeps (\d+)ms before responding (\d+)$`, func(ms, status int) error { return f.aTestHTTPServerThatSleepsMsThenResponds(ms, status) })
	ctx.Step(`^I run the CLI with the server URL, requests (\d+) and concurrency (\d+)$`, func(r, c int) error { return f.iRunTheCLIWithServerURLRequestsAndConcurrency(r, c) })
	ctx.Step(`^the report shows total requests (\d+)$`, func(n int) error { return f.theReportShowsTotalRequests(n) })
	ctx.Step(`^the report shows 200 OK: (\d+)$`, func(n int) error { return f.theReportShows200OK(n) })
	ctx.Step(`^the report shows (\d+): (\d+)$`, func(code, n int) error { return f.theReportShowsStatusCount(code, n) })
	ctx.Step(`^the report total time is at least (\d+)ms$`, func(n int) error { return f.theReportTotalAtLeastMs(n) })
	// cleanup will be done after running the CLI in the step
}
