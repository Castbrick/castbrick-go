package castbrick

import (
	"context"
	"net/url"
	"strconv"
)

type WebhooksResource struct {
	c *Client
}

func (r *WebhooksResource) List(ctx context.Context, page, pageSize int) (*PagedResult[Webhook], error) {
	q := url.Values{}
	if page > 0 {
		q.Set("pageNumber", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}

	var res PagedResult[Webhook]
	if err := r.c.getQuery(ctx, "/webhooks", q, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *WebhooksResource) Get(ctx context.Context, id string) (*Webhook, error) {
	path, err := pathf("/webhooks/%s", id)
	if err != nil {
		return nil, err
	}
	var res Webhook
	if err := r.c.get(ctx, path, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *WebhooksResource) Create(ctx context.Context, data CreateWebhookRequest) (string, error) {
	var res string
	if err := r.c.post(ctx, "/webhooks", data, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *WebhooksResource) Toggle(ctx context.Context, id string) error {
	path, err := pathf("/webhooks/%s/toggle", id)
	if err != nil {
		return err
	}
	return r.c.put(ctx, path, map[string]any{}, nil)
}

func (r *WebhooksResource) Test(ctx context.Context, webhookID, payload string) error {
	data := map[string]string{
		"webhookId": webhookID,
		"payload":   payload,
	}
	return r.c.post(ctx, "/webhooks/test", data, nil)
}

func (r *WebhooksResource) ListLogs(ctx context.Context, id string, limit int) ([]WebhookLog, error) {
	path, err := pathf("/webhooks/%s/logs", id)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}

	var res []WebhookLog
	if err := r.c.getQuery(ctx, path, q, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *WebhooksResource) Retry(ctx context.Context, logID string) error {
	path, err := pathf("/webhooks/logs/%s/retry", logID)
	if err != nil {
		return err
	}
	return r.c.post(ctx, path, map[string]any{}, nil)
}

func (r *WebhooksResource) Delete(ctx context.Context, id string) error {
	path, err := pathf("/webhooks/%s", id)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
