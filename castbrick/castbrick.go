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
func New(apiKey string) *CastBrick {
	return NewWithOptions(apiKey, defaultBaseURL, nil)
}

// NewWithOptions creates a CastBrick client with a custom base URL and HTTP client.
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
