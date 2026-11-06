package castbrick

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

// BroadcastsResource handles broadcast operations.
type BroadcastsResource struct{ c *Client }

type CreateBroadcastOptions struct {
	Name          string
	Message       string
	ContactListID string
	SenderID      string
}

type UpdateBroadcastOptions struct {
	Name          string
	Message       string
	ContactListID string
	SenderID      string
	ScheduleAt    *time.Time
}

func (r *BroadcastsResource) List(ctx context.Context, page, pageSize int) (*PagedResult[Broadcast], error) {
	q := url.Values{}
	q.Set("pageNumber", strconv.Itoa(page))
	q.Set("pageSize", strconv.Itoa(pageSize))
	var out PagedResult[Broadcast]
	if err := r.c.getQuery(ctx, "/broadcasts", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *BroadcastsResource) Get(ctx context.Context, id string) (*Broadcast, error) {
	path, err := pathf("/broadcasts/%s", id)
	if err != nil {
		return nil, err
	}
	var out Broadcast
	if err := r.c.get(ctx, path, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *BroadcastsResource) Create(ctx context.Context, opts CreateBroadcastOptions) (string, error) {
	body := map[string]any{
		"name":    opts.Name,
		"message": opts.Message,
	}
	if opts.ContactListID != "" {
		body["contactListId"] = opts.ContactListID
	}
	if opts.SenderID != "" {
		body["senderId"] = opts.SenderID
	}
	var out string
	if err := r.c.post(ctx, "/broadcasts", body, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (r *BroadcastsResource) Update(ctx context.Context, id string, opts UpdateBroadcastOptions) (string, error) {
	path, err := pathf("/broadcasts/%s", id)
	if err != nil {
		return "", err
	}
	body := map[string]any{
		"name":    opts.Name,
		"message": opts.Message,
	}
	if opts.ContactListID != "" {
		body["contactListId"] = opts.ContactListID
	}
	if opts.SenderID != "" {
		body["senderId"] = opts.SenderID
	}
	if opts.ScheduleAt != nil {
		body["scheduleAt"] = opts.ScheduleAt.UTC().Format(time.RFC3339)
	}
	var out string
	if err := r.c.put(ctx, path, body, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (r *BroadcastsResource) Send(ctx context.Context, id string) error {
	path, err := pathf("/broadcasts/%s/send", id)
	if err != nil {
		return err
	}
	return r.c.post(ctx, path, map[string]any{}, nil)
}

func (r *BroadcastsResource) Cancel(ctx context.Context, id string) error {
	path, err := pathf("/broadcasts/%s/cancel", id)
	if err != nil {
		return err
	}
	return r.c.post(ctx, path, map[string]any{}, nil)
}

func (r *BroadcastsResource) Duplicate(ctx context.Context, id string) (string, error) {
	path, err := pathf("/broadcasts/%s/duplicate", id)
	if err != nil {
		return "", err
	}
	var out string
	if err := r.c.post(ctx, path, map[string]any{}, &out); err != nil {
		return "", err
	}
	return out, nil
}

func (r *BroadcastsResource) Delete(ctx context.Context, id string) error {
	path, err := pathf("/broadcasts/%s", id)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
