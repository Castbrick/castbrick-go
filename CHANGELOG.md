# Changelog

## 0.3.1 — Unreleased

### Bug Fixes
- **Build** — `v0.3.0` was tagged from a commit where `billing.go`, `segments.go`, `templates.go` and `webhooks.go` called `Client.newRequest` and `Client.do(req, out)`, which did not exist yet. The module did not compile for any consumer. `v0.3.0` is now retracted in `go.mod`.
- **Results on error** — methods now return a `nil` result together with the error, instead of a pointer to an empty value.

### Security
- **Path injection** — resource IDs are validated and escaped with `url.PathEscape` before being added to the request path. Values such as `../webhooks/abc` or `id?status=all` used to reach the API as a different endpoint or query. They are now rejected with `ErrInvalidID`.
- **Unbounded response reads** — response bodies are capped at 10 MiB (`ErrResponseTooLarge`). `APIError.Body` is capped at 64 KiB, and `APIError.Error()` at 512 bytes with control characters removed.
- **No default timeout** — the default HTTP client now has a 30 second timeout.
- **API key forwarded on redirects** — the default HTTP client refuses redirects that change scheme, host or port, so the `Authorization` header is never sent to another origin.
- **Insecure base URL** — the base URL must use `https` (plain `http` is only allowed for loopback hosts) and must not contain credentials, a query or a fragment. An empty API key is rejected. Invalid settings make every call fail with `ErrInvalidConfig`.

### Improvements
- `APIError.Message` — the error message parsed from the JSON body (`message`, `error`, `detail` or `title`).
- A `User-Agent: castbrick-go` header is sent with every request.
- Query strings are built with `url.Values` in every resource.
- CI runs on every push and pull request: `gofmt`, `go build`, `go vet` and `go test -race` on the `go.mod` version and on stable Go, plus `govulncheck`. GitHub Actions are pinned to commit SHAs and use read-only permissions.

### Migration notes
- No public signatures changed.
- Code that read fields from the result of a failed call must check `err` first. The result is now `nil`.
- Code that used `http://` against a non-local host, or an empty API key, now gets `ErrInvalidConfig`.

---

## 0.2.0 — 2026-07-18

### New Features
- **Billing API (`BillingResource`)** — access credit balance, list packages, initiate payments, and view transactions.
- **Segments API (`SegmentsResource`)** — list, create, update, and delete dynamic audience segments.
- **Templates API (`TemplatesResource`)** — list, create, update, and delete message templates.
- **Webhooks API (`WebhooksResource`)** — register endpoints, toggle status, test delivery, list logs, and retry failed webhooks.
- **`models.go`** — added structs required for `Template`, `Webhook`, `Segment`, and `Billing` payloads and responses.

---

## 0.1.3 — 2026-06-01

### Bug Fixes
- **`SMS.CancelScheduled()`** — fixed endpoint from `POST /sms/cancel-scheduled` to `DELETE /sms/{id}` (was completely broken)
- **`Contacts.CreateList()`** — now correctly returns `(string, error)` (the new list ID) instead of trying to deserialize a `*ContactList` object
- **`Contacts.Create()`** — removed unsupported `emails` parameter; API only accepts `phoneNumbers`. Signature changed from `Create(ctx, emails, phoneNumbers string)` to `Create(ctx, phoneNumbers string)`
- **`defaultBaseURL`** — corrected to `https://api.castbrick.co` (was `https://api.castbrick.co/v1`, causing all requests to 404)

### New Features
- **`SMS.Send()`** — added `Fallback *bool` field to `SendSmsOptions` to control sender ID fallback behaviour
- **`SMS.List()`** — signature changed to accept `ListSmsOptions` struct with optional `Status`, `Phone`, `From` and `To` filter fields
- **`ListSmsOptions`** — new exported type for `SMS.List()` parameters

### Breaking Changes
- `SMS.List(ctx, page, pageSize int)` → `SMS.List(ctx, ListSmsOptions{Page, PageSize, ...})`
- `Contacts.Create(ctx, emails, phoneNumbers string)` → `Contacts.Create(ctx, phoneNumbers string)`
- `Contacts.CreateList(ctx, name string) (*ContactList, error)` → `Contacts.CreateList(ctx, name string) (string, error)`

---

## 0.1.2

- Initial release.
