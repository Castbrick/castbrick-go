# castbrick-go

Official Go SDK for the [CastBrick](https://castbrick.co) API — send SMS, manage contacts and run broadcasts.

## Installation

```bash
go get github.com/IldySilva/castbrick-go
```

## Quick start

```go
import "github.com/IldySilva/castbrick-go/castbrick"

cb := castbrick.New("your_api_key")

result, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{
    To:      []string{"+244923000000"},
    Content: "Hello from CastBrick!",
})
```

## SMS

```go
// Send
fallback := true
result, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{
    To:       []string{"+244923000000"},
    Content:  "Your OTP is 1234",
    SenderID: "MyApp",     // optional — your approved Sender ID
    Fallback: &fallback,   // optional — fall back to CastBrick sender if SenderID unavailable
})

// List (simple)
page, err := cb.SMS.List(ctx, castbrick.ListSmsOptions{Page: 1, PageSize: 20})
fmt.Println(page.TotalCount)

// List with filters
from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
page, err = cb.SMS.List(ctx, castbrick.ListSmsOptions{
    Page:     1,
    PageSize: 20,
    Status:   "delivered",        // pending | sent | delivered | failed | scheduled
    Phone:    "+244923000000",
    From:     &from,
})

// Cancel a scheduled SMS
err = cb.SMS.CancelScheduled(ctx, "message-id")
```

## Contacts

```go
// List (with optional search)
page, err := cb.Contacts.List(ctx, 1, 20, "john")

// Get
contact, err := cb.Contacts.Get(ctx, "contact-id")

// Create — comma or newline-separated phone numbers
err = cb.Contacts.Create(ctx, "+244923000000,+244912000000")

// Delete
err = cb.Contacts.Delete(ctx, "contact-id")

// Contact lists — CreateList returns the new list ID (string)
listID, err := cb.Contacts.CreateList(ctx, "VIP Customers")
err = cb.Contacts.AddToList(ctx, listID, contact.ID)
err = cb.Contacts.RemoveFromList(ctx, listID, contact.ID)
```

## Broadcasts

```go
// Create and send immediately
id, err := cb.Broadcasts.Create(ctx, castbrick.CreateBroadcastOptions{
    Name:          "Black Friday",
    Message:       "50% off everything today!",
    ContactListID: "list-id", // optional
    SenderID:      "MyApp",   // optional
})
err = cb.Broadcasts.Send(ctx, id)

// Schedule for later
t := time.Date(2026, 11, 28, 9, 0, 0, 0, time.UTC)
_, err = cb.Broadcasts.Update(ctx, id, castbrick.UpdateBroadcastOptions{
    Name:       "Black Friday",
    Message:    "50% off everything today!",
    ScheduleAt: &t,
})
err = cb.Broadcasts.Send(ctx, id)

// Other operations
err = cb.Broadcasts.Cancel(ctx, id)
newID, err := cb.Broadcasts.Duplicate(ctx, id)
err = cb.Broadcasts.Delete(ctx, id)
```

## Billing

```go
// Get current credit balance
balance, err := cb.Billing.GetBalance(ctx)

// List available credit packages
packages, err := cb.Billing.ListPackages(ctx)

// List credit consumption/transaction history
history, err := cb.Billing.ListTransactions(ctx, "")
```

## Segments, Templates & Webhooks

The SDK also provides native access to `cb.Segments`, `cb.Templates`, and `cb.Webhooks`.

```go
// Example: Create a Webhook
err = cb.Webhooks.Create(ctx, castbrick.CreateWebhookOptions{
    Url:    "https://yourapp.com/webhooks/castbrick",
    Events: []string{"sms.delivered", "sms.failed"},
})
```

## Error handling

```go
_, err := cb.SMS.Send(ctx, castbrick.SendSmsOptions{...})

var apiErr *castbrick.APIError
switch {
case errors.As(err, &apiErr):
    // apiErr.Message holds the API error message, apiErr.Body the raw body (max 64 KiB).
    fmt.Printf("%d: %s\n", apiErr.StatusCode, apiErr.Message)
    // 401 → invalid or revoked API key
    // 402 → insufficient credits
    // 422 → validation error
case errors.Is(err, castbrick.ErrInvalidID):
    // An ID was empty or contained characters such as "/", "?", "#" or "..".
case errors.Is(err, castbrick.ErrInvalidConfig):
    // Empty API key or an unsafe base URL (see "Security").
}
```

On error, methods return a `nil` result together with the error.

## Custom HTTP client

```go
// An empty base URL selects the production API.
cb := castbrick.NewWithOptions("your_api_key", "", &http.Client{
    Timeout: 10 * time.Second,
})
```

A custom client is used as-is: always set a `Timeout`.

## Security

The SDK applies these safeguards by default:

- **Resource IDs are validated and escaped.** IDs that are empty, longer than 256 characters, or contain `/`, `\`, `?`, `#`, `%`, control characters, `.` or `..` are rejected with `ErrInvalidID` before any request is sent. This means an ID taken from user input cannot redirect a call to another endpoint.
- **HTTPS is required.** The base URL must use `https`. Plain `http` is only accepted for `localhost`, `127.0.0.1` and `::1`. URLs with credentials, a query or a fragment are rejected.
- **The API key stays on its origin.** The default HTTP client refuses redirects that change scheme or host (including the port).
- **Requests are bounded.** The default client times out after 30 seconds, and responses over 10 MiB are rejected with `ErrResponseTooLarge`.
- **Errors are safe to log.** `APIError.Error()` is truncated to 512 bytes and has control characters removed.

Recommendations for applications that use the SDK:

- Keep the API key in a secret manager or environment variable, never in source code.
- Validate the format of IDs that come from your own users (for example as UUIDs) before passing them to the SDK, and check that they belong to the current user.
- Never fill `InitiatePaymentRequest` redirect URLs from user input without an allowlist.
- Pass a `context` with a deadline to every call.
- Build your application with a supported, patched Go release, because security fixes in `net/http` and `crypto/tls` come from your toolchain.

To report a vulnerability, see [SECURITY.md](SECURITY.md).

## License

MIT
