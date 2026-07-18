# castbrick-go

Official Go SDK for the [CastBrick](https://castbrick.com) API — send SMS, manage contacts and run broadcasts.

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
if errors.As(err, &apiErr) {
    fmt.Printf("%d: %s\n", apiErr.StatusCode, apiErr.Body)
    // 401 → invalid or revoked API key
    // 402 → insufficient credits
    // 422 → validation error
}
```

## Custom HTTP client

```go
cb := castbrick.NewWithOptions("your_api_key", "https://api.castbrick.co", &http.Client{
    Timeout: 10 * time.Second,
})
```

## License

MIT
