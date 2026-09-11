package castbrick_test

import (
	"context"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/IldySilva/castbrick-go/castbrick"
)

func TestMain(m *testing.M) {
	baseURL := os.Getenv("CASTBRICK_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8787/v1"
	}
	healthURL := strings.TrimRight(strings.TrimSuffix(baseURL, "/v1"), "/") + "/health"
	client := http.Client{Timeout: 500 * time.Millisecond}
	resp, err := client.Get(healthURL)
	var cmd *exec.Cmd
	if err != nil || resp.StatusCode != 200 {
		serverScript := "../tools/mock-server/server.js"
		if _, statErr := os.Stat(serverScript); statErr == nil {
			cmd = exec.Command("node", serverScript, "8787")
			_ = cmd.Start()
			time.Sleep(500 * time.Millisecond)
		}
	} else if resp != nil {
		_ = resp.Body.Close()
	}

	code := m.Run()

	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
	os.Exit(code)
}

func getClient(t *testing.T) *castbrick.CastBrick {
	t.Helper()
	baseURL := os.Getenv("CASTBRICK_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8787/v1"
	}
	apiKey := os.Getenv("CASTBRICK_API_KEY")
	if apiKey == "" {
		apiKey = "test_api_key"
	}
	return castbrick.NewWithOptions(apiKey, baseURL, nil)
}

func TestSMS_Send(t *testing.T) {
	cb := getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{
		To:      []string{"+244923000000"},
		Content: "Hello from Go SDK Test!",
	})
	if err != nil {
		t.Fatalf("SMS.Send failed: %v", err)
	}

	if res.MessageID == "" {
		t.Errorf("expected non-empty message ID")
	}
	if res.Status != "Queued" {
		t.Errorf("expected status 'Queued', got %s", res.Status)
	}
}

func TestSMS_List(t *testing.T) {
	cb := getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cb.SMS.List(ctx, castbrick.ListSmsOptions{Page: 1, PageSize: 10})
	if err != nil {
		t.Fatalf("SMS.List failed: %v", err)
	}
	if res.TotalCount < 1 {
		t.Errorf("expected totalCount >= 1, got %d", res.TotalCount)
	}
}

func TestContacts_List(t *testing.T) {
	cb := getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cb.Contacts.List(ctx, 1, 10, "")
	if err != nil {
		t.Fatalf("Contacts.List failed: %v", err)
	}
	if len(res.Items) == 0 {
		t.Errorf("expected at least 1 contact")
	}
}

func TestBroadcasts_List(t *testing.T) {
	cb := getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cb.Broadcasts.List(ctx, 1, 10)
	if err != nil {
		t.Fatalf("Broadcasts.List failed: %v", err)
	}
	if len(res.Items) == 0 {
		t.Errorf("expected at least 1 broadcast")
	}
}

func TestBilling_GetBalance(t *testing.T) {
	cb := getClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := cb.Billing.GetBalance(ctx)
	if err != nil {
		t.Fatalf("Billing.GetBalance failed: %v", err)
	}
	if res.Balance != 50000.0 {
		t.Errorf("expected balance 50000.0, got %f", res.Balance)
	}
}

func TestUnauthorized(t *testing.T) {
	baseURL := os.Getenv("CASTBRICK_BASE_URL")
	if baseURL == "" {
		baseURL = "http://127.0.0.1:8787/v1"
	}
	cb := castbrick.NewWithOptions("invalid_key", baseURL, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{
		To:      []string{"+244923000000"},
		Content: "fail",
	})
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}

	apiErr, ok := err.(*castbrick.APIError)
	if !ok {
		t.Fatalf("expected *castbrick.APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected status code 401, got %d", apiErr.StatusCode)
	}
}
