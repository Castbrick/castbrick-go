# Changelog

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
