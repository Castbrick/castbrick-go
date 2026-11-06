package castbrick

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// SmsResource handles SMS operations.
type SmsResource struct{ c *Client }

type SendSmsOptions struct {
	To            []string
	Content       string
	SenderID      string
	ScheduledAt   *time.Time
	ContactListID string
	// Fallback controls whether to fall back to the CastBrick default sender
	// when SenderID is unavailable. Defaults to true on the server.
	Fallback *bool
}

// ListSmsOptions contains optional filters for listing SMS messages.
type ListSmsOptions struct {
	Page     int
	PageSize int
	Status   string
	Phone    string
	From     *time.Time
	To       *time.Time
}

func (r *SmsResource) Send(ctx context.Context, opts SendSmsOptions) (*SendSmsResponse, error) {
	body := map[string]any{
		"recipients": opts.To,
		"content":    opts.Content,
	}
	if opts.SenderID != "" {
		body["senderId"] = opts.SenderID
	}
	if opts.ScheduledAt != nil {
		body["scheduledAt"] = opts.ScheduledAt.UTC().Format(time.RFC3339)
	}
	if opts.ContactListID != "" {
		body["contactListId"] = opts.ContactListID
	}
	if opts.Fallback != nil {
		body["fallback"] = *opts.Fallback
	}

	var out SendSmsResponse
	if err := r.c.post(ctx, "/sms/send", body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *SmsResource) List(ctx context.Context, opts ListSmsOptions) (*PagedResult[SmsMessage], error) {
	page := opts.Page
	if page < 1 {
		page = 1
	}
	pageSize := opts.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	q := url.Values{}
	q.Set("pageNumber", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	if opts.Status != "" {
		q.Set("status", opts.Status)
	}
	if opts.Phone != "" {
		q.Set("phone", opts.Phone)
	}
	if opts.From != nil {
		q.Set("from", opts.From.UTC().Format(time.RFC3339))
	}
	if opts.To != nil {
		q.Set("to", opts.To.UTC().Format(time.RFC3339))
	}

	var out PagedResult[SmsMessage]
	if err := r.c.getQuery(ctx, "/sms", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CancelScheduled cancels a scheduled SMS by its message ID.
func (r *SmsResource) CancelScheduled(ctx context.Context, messageID string) error {
	path, err := pathf("/sms/%s", messageID)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
