package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

// TestConcurrentRequestCount verifies that the request counter is safe under
// concurrent access. Run with -race to detect the data race.
func TestConcurrentRequestCount(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Get(srv.URL + "/healthz")
			if err != nil {
				t.Errorf("GET /healthz: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()

	if got := api.RequestCount(); got != goroutines {
		t.Errorf("RequestCount() = %d, want %d", got, goroutines)
	}
}

// TestConcurrentValidateCount ensures validate endpoint also increments safely.
func TestConcurrentValidateCount(t *testing.T) {
	api := New()
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			resp, err := http.Post(srv.URL+"/api/validate", "application/json",
				strings.NewReader(`{"expr":"* * * * *"}`))
			if err != nil {
				t.Errorf("POST: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}
	wg.Wait()

	if got := api.RequestCount(); got != goroutines {
		t.Errorf("RequestCount() = %d, want %d", got, goroutines)
	}
}
