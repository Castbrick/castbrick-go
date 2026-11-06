package castbrick_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/IldySilva/castbrick-go/castbrick"
)

// recorder is a test server that records every request it receives.
type recorder struct {
	srv   *httptest.Server
	hits  atomic.Int32
	paths []string
}

func newRecorder(t *testing.T, status int, body string) *recorder {
	t.Helper()
	rec := &recorder{}
	rec.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec.hits.Add(1)
		rec.paths = append(rec.paths, r.Method+" "+r.URL.EscapedPath()+"?"+r.URL.RawQuery)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(rec.srv.Close)
	return rec
}

func TestMaliciousIDsAreRejectedBeforeAnyRequest(t *testing.T) {
	rec := newRecorder(t, http.StatusOK, `{}`)
	cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)
	ctx := context.Background()

	malicious := []string{
		"",
		".",
		"..",
		"../webhooks/abc",
		"abc/../../billing/balance",
		"abc?status=all",
		"abc#fragment",
		"..%2Fwebhooks",
		`..\webhooks`,
		"abc\r\nX-Injected: 1",
		strings.Repeat("a", 257),
	}

	calls := map[string]func(id string) error{
		"Broadcasts.Get":          func(id string) error { _, err := cb.Broadcasts.Get(ctx, id); return err },
		"Broadcasts.Delete":       func(id string) error { return cb.Broadcasts.Delete(ctx, id) },
		"Broadcasts.Send":         func(id string) error { return cb.Broadcasts.Send(ctx, id) },
		"Contacts.Get":            func(id string) error { _, err := cb.Contacts.Get(ctx, id); return err },
		"Contacts.Delete":         func(id string) error { return cb.Contacts.Delete(ctx, id) },
		"Contacts.RemoveFromList": func(id string) error { return cb.Contacts.RemoveFromList(ctx, "list-1", id) },
		"Contacts.AddToList":      func(id string) error { return cb.Contacts.AddToList(ctx, id, "contact-1") },
		"SMS.CancelScheduled":     func(id string) error { return cb.SMS.CancelScheduled(ctx, id) },
		"Segments.Delete":         func(id string) error { return cb.Segments.Delete(ctx, id) },
		"Templates.Delete":        func(id string) error { return cb.Templates.Delete(ctx, id) },
		"Webhooks.Delete":         func(id string) error { return cb.Webhooks.Delete(ctx, id) },
		"Webhooks.Retry":          func(id string) error { return cb.Webhooks.Retry(ctx, id) },
		"Webhooks.ListLogs":       func(id string) error { _, err := cb.Webhooks.ListLogs(ctx, id, 10); return err },
	}

	for name, call := range calls {
		for _, id := range malicious {
			if err := call(id); !errors.Is(err, castbrick.ErrInvalidID) {
				t.Errorf("%s(%q): expected ErrInvalidID, got %v", name, id, err)
			}
		}
	}
	if n := rec.hits.Load(); n != 0 {
		t.Fatalf("expected no request to reach the server, got %d: %v", n, rec.paths)
	}
}

func TestValidIDsAreEscapedIntoASingleSegment(t *testing.T) {
	rec := newRecorder(t, http.StatusNoContent, ``)
	cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)

	for _, id := range []string{"3f2c1a9e-0b7d-4c1e-9a55-2d6f0e8b1c34", "msg_123", "a b"} {
		if err := cb.SMS.CancelScheduled(context.Background(), id); err != nil {
			t.Fatalf("CancelScheduled(%q): %v", id, err)
		}
	}
	want := []string{
		"DELETE /sms/3f2c1a9e-0b7d-4c1e-9a55-2d6f0e8b1c34?",
		"DELETE /sms/msg_123?",
		"DELETE /sms/a%20b?",
	}
	for i, w := range want {
		if rec.paths[i] != w {
			t.Errorf("request %d: got %q, want %q", i, rec.paths[i], w)
		}
	}
}

func TestQueryParametersAreEncoded(t *testing.T) {
	rec := newRecorder(t, http.StatusOK, `{"items":[]}`)
	cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)

	if _, err := cb.Contacts.List(context.Background(), 1, 10, "a&pageSize=100000"); err != nil {
		t.Fatal(err)
	}
	got := rec.paths[0]
	if !strings.Contains(got, "search=a%26pageSize%3D100000") || strings.Count(got, "pageSize=") != 1 {
		t.Fatalf("search value was not encoded safely: %s", got)
	}
}

func TestInsecureConfigurationIsRejected(t *testing.T) {
	cases := map[string]struct{ key, baseURL string }{
		"plain http to remote host": {"key", "http://api.castbrick.co/v1"},
		"credentials in URL":        {"key", "https://user:pass@api.castbrick.co"},
		"query in URL":              {"key", "https://api.castbrick.co?x=1"},
		"unsupported scheme":        {"key", "ftp://api.castbrick.co"},
		"missing host":              {"key", "https://"},
		"empty API key":             {"", "https://api.castbrick.co/v1"},
		"blank API key":             {"   ", "https://api.castbrick.co/v1"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cb := castbrick.NewWithOptions(tc.key, tc.baseURL, nil)
			_, err := cb.Billing.GetBalance(context.Background())
			if !errors.Is(err, castbrick.ErrInvalidConfig) {
				t.Fatalf("expected ErrInvalidConfig, got %v", err)
			}
		})
	}
}

func TestLoopbackHTTPIsAllowed(t *testing.T) {
	rec := newRecorder(t, http.StatusOK, `{"balance":1}`)
	// httptest servers listen on http://127.0.0.1; a trailing slash must not break paths.
	cb := castbrick.NewWithOptions("key", rec.srv.URL+"/", nil)
	if _, err := cb.Billing.GetBalance(context.Background()); err != nil {
		t.Fatalf("loopback http should be allowed: %v", err)
	}
	if rec.paths[0] != "GET /billing/balance?" {
		t.Fatalf("unexpected path %q", rec.paths[0])
	}
}

func TestOversizedResponseIsRejected(t *testing.T) {
	rec := newRecorder(t, http.StatusOK, `"`+strings.Repeat("A", 11<<20)+`"`)
	cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)
	_, err := cb.Billing.GetBalance(context.Background())
	if !errors.Is(err, castbrick.ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got %v", err)
	}
}

func TestAPIErrorIsBoundedAndSafeToLog(t *testing.T) {
	t.Run("raw body is truncated", func(t *testing.T) {
		rec := newRecorder(t, http.StatusInternalServerError, strings.Repeat("A", 5<<20))
		cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)
		_, err := cb.Billing.GetBalance(context.Background())

		var apiErr *castbrick.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T", err)
		}
		if len(apiErr.Body) > 64<<10 {
			t.Errorf("APIError.Body not truncated: %d bytes", len(apiErr.Body))
		}
		if len(err.Error()) > 600 {
			t.Errorf("Error() not truncated: %d bytes", len(err.Error()))
		}
	})

	t.Run("message is extracted and control characters are removed", func(t *testing.T) {
		rec := newRecorder(t, http.StatusBadRequest, `{"message":"bad\nFAKE LOG LINE"}`)
		cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)
		_, err := cb.Billing.GetBalance(context.Background())

		var apiErr *castbrick.APIError
		if !errors.As(err, &apiErr) {
			t.Fatalf("expected *APIError, got %T", err)
		}
		if apiErr.StatusCode != http.StatusBadRequest || apiErr.Message != "bad\nFAKE LOG LINE" {
			t.Errorf("unexpected APIError: %+v", apiErr)
		}
		if strings.ContainsAny(err.Error(), "\r\n") {
			t.Errorf("Error() contains a line break: %q", err.Error())
		}
	})
}

func TestFailedCallsReturnNilResult(t *testing.T) {
	rec := newRecorder(t, http.StatusUnauthorized, `{"message":"unauthorized"}`)
	cb := castbrick.NewWithOptions("key", rec.srv.URL, nil)
	ctx := context.Background()

	if res, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{To: []string{"+244923000000"}, Content: "x"}); err == nil || res != nil {
		t.Errorf("SMS.Send: expected (nil, error), got (%v, %v)", res, err)
	}
	if res, err := cb.Broadcasts.Get(ctx, "b-1"); err == nil || res != nil {
		t.Errorf("Broadcasts.Get: expected (nil, error), got (%v, %v)", res, err)
	}
	if res, err := cb.Contacts.List(ctx, 1, 10, ""); err == nil || res != nil {
		t.Errorf("Contacts.List: expected (nil, error), got (%v, %v)", res, err)
	}
}

func TestDefaultClientDoesNotLeakAPIKeyOnCrossOriginRedirect(t *testing.T) {
	var leaked atomic.Bool
	evil := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			leaked.Store(true)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	defer evil.Close()

	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Same IP, different port: a different origin that Go's own redirect
		// policy would still treat as the same host for header forwarding.
		http.Redirect(w, r, evil.URL+"/collect", http.StatusTemporaryRedirect)
	}))
	defer api.Close()

	cb := castbrick.NewWithOptions("secret-key", api.URL, nil)
	if _, err := cb.Billing.GetBalance(context.Background()); err == nil {
		t.Fatal("expected the cross-origin redirect to be refused")
	}
	if leaked.Load() {
		t.Fatal("API key was forwarded to another origin")
	}
}

func TestUserAgentIsSent(t *testing.T) {
	var ua atomic.Value
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua.Store(r.Header.Get("User-Agent"))
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	if _, err := castbrick.NewWithOptions("key", srv.URL, nil).Billing.GetBalance(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got, _ := ua.Load().(string); !strings.HasPrefix(got, "castbrick-go") {
		t.Fatalf("unexpected User-Agent %q", got)
	}
}
