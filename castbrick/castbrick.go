// Package castbrick is the official Go SDK for the CastBrick API.
//
// Usage:
//
//	cb := castbrick.New("your_api_key")
//
//	result, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{
//	    To:      []string{"+244923000000"},
//	    Content: "Hello from CastBrick!",
//	})
package castbrick

import "net/http"

// CastBrick is the main SDK client.
type CastBrick struct {
	SMS        *SmsResource
	Contacts   *ContactsResource
	Broadcasts *BroadcastsResource
	Billing    *BillingResource
	Segments   *SegmentsResource
	Templates  *TemplatesResource
	Webhooks   *WebhooksResource
}

// New creates a CastBrick client with the default HTTP client.
//
// The default HTTP client has a 30 second timeout and refuses redirects that
// change host or scheme, so the API key is never forwarded to another origin.
func New(apiKey string) *CastBrick {
	return NewWithOptions(apiKey, defaultBaseURL, nil)
}

// NewWithOptions creates a CastBrick client with a custom base URL and HTTP client.
//
// An empty baseURL selects the production API. The base URL must use https;
// plain http is only accepted for loopback hosts (localhost, 127.0.0.1, ::1).
// If httpClient is nil, the SDK default client is used; a custom client is used
// as-is, so give it a Timeout.
//
// An invalid API key or base URL does not panic: every call made with the
// returned client fails with an error wrapping ErrInvalidConfig.
func NewWithOptions(apiKey, baseURL string, httpClient *http.Client) *CastBrick {
	c := newClient(apiKey, baseURL, httpClient)
	return &CastBrick{
		SMS:        &SmsResource{c: c},
		Contacts:   &ContactsResource{c: c},
		Broadcasts: &BroadcastsResource{c: c},
		Billing:    &BillingResource{c: c},
		Segments:   &SegmentsResource{c: c},
		Templates:  &TemplatesResource{c: c},
		Webhooks:   &WebhooksResource{c: c},
	}
}
