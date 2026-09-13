package loadtest

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

func runAndOutput(args []string) (int, string) {
	var buf bytes.Buffer
	code := Run(args, &buf)
	return code, buf.String()
}

func TestRun_MissingURL(t *testing.T) {
	code, out := runAndOutput([]string{})
	if code == 0 {
		t.Fatalf("expected non-zero exit code when url missing")
	}
	if !strings.Contains(out, "--url é obrigatório") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRun_InvalidParams(t *testing.T) {
	code, out := runAndOutput([]string{"--url=http://example.local", "--requests=0", "--concurrency=1"})
	if code == 0 {
		t.Fatalf("expected non-zero exit code when requests=0")
	}
	if !strings.Contains(out, "--requests e --concurrency devem ser maiores que 0") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestRun_ConcurrencyAdjusted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	code, out := runAndOutput([]string{"--url=" + srv.URL, "--requests=2", "--concurrency=5"})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d output:%s", code, out)
	}
	if !strings.Contains(out, "Total de requests: 2") {
		t.Fatalf("unexpected total: %s", out)
	}
	if !strings.Contains(out, "200 OK: 2") {
		t.Fatalf("unexpected 200 count: %s", out)
	}
}

func TestRun_MixedStatuses(t *testing.T) {
	// cycle through 200,404,500
	statuses := []int{200, 404, 500}
	var cnt int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := atomic.AddInt64(&cnt, 1) - 1
		code := statuses[int(i)%len(statuses)]
		w.WriteHeader(code)
		_, _ = w.Write([]byte(strconv.Itoa(code)))
	}))
	defer srv.Close()

	code, out := runAndOutput([]string{"--url=" + srv.URL, "--requests=9", "--concurrency=3"})
	if code != 0 {
		t.Fatalf("expected zero exit code, got %d output:%s", code, out)
	}
	for _, s := range statuses {
		want := "" + strconv.Itoa(s) + ": 3"
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in output, got:\n%s", want, out)
		}
	}
}
