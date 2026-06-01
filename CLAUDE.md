# castbrick-go

Official Go SDK for the CastBrick API. Module path: `github.com/IldySilva/castbrick-go`.

## Commands

```bash
go build ./...    # compile
go test ./...     # run tests
go vet ./...      # lint
```

## Version

Versioned via git tags: `v0.1.3`. No version field in `go.mod`.
```bash
git tag v0.1.x && git push origin v0.1.x
```

## Source files

```
castbrick/
├── castbrick.go    # CastBrick struct + New() + NewWithOptions()
├── client.go       # Client, do(), get(), post(), put(), delete(), APIError
├── models.go       # PagedResult, SendSmsResponse, SmsMessage, Contact, ContactList, Broadcast
├── sms.go          # SmsResource: Send, List, CancelScheduled
├── contacts.go     # ContactsResource: List, Get, Create, Delete, CreateList, AddToList, RemoveFromList
└── broadcasts.go   # BroadcastsResource: List, Get, Create, Update, Send, Cancel, Duplicate, Delete
```

## API base URL

`defaultBaseURL = "https://api.castbrick.co"` — **no `/v1`**.

## Correct API facts

- `SMS.CancelScheduled(ctx, messageID)` → `DELETE /sms/{id}`
- `Contacts.CreateList(ctx, name)` → returns `(string, error)` (ID only, not `*ContactList`)
- `Contacts.Create(ctx, phoneNumbers string)` — single param, no `emails`
- `SMS.List(ctx, ListSmsOptions{Status, Phone, From, To, Page, PageSize})` — use options struct
- `SendSmsOptions.Fallback *bool` — optional fallback control

## Breaking changes from v0.1.2

- `SMS.List(ctx, page, pageSize int)` → `SMS.List(ctx, ListSmsOptions{...})`
- `Contacts.Create(ctx, emails, phoneNumbers string)` → `Contacts.Create(ctx, phoneNumbers string)`
- `Contacts.CreateList(ctx, name) (*ContactList, error)` → `(string, error)`
