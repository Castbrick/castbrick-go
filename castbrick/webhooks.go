package castbrick

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type WebhooksResource struct {
	c *Client
}

func (r *WebhooksResource) List(ctx context.Context, page, pageSize int) (*PagedResult[Webhook], error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/webhooks", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if page > 0 {
		q.Add("pageNumber", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Add("pageSize", strconv.Itoa(pageSize))
	}
	req.URL.RawQuery = q.Encode()

	var res PagedResult[Webhook]
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *WebhooksResource) Get(ctx context.Context, id string) (*Webhook, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/webhooks/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var res Webhook
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *WebhooksResource) Create(ctx context.Context, data CreateWebhookRequest) (string, error) {
	req, err := r.c.newRequest(ctx, http.MethodPost, "/webhooks", data)
	if err != nil {
		return "", err
	}

	var res string
	if err := r.c.do(req, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *WebhooksResource) Toggle(ctx context.Context, id string) error {
	req, err := r.c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/webhooks/%s/toggle", id), map[string]interface{}{})
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}

func (r *WebhooksResource) Test(ctx context.Context, webhookID, payload string) error {
	data := map[string]string{
		"webhookId": webhookID,
		"payload":   payload,
	}
	req, err := r.c.newRequest(ctx, http.MethodPost, "/webhooks/test", data)
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}

func (r *WebhooksResource) ListLogs(ctx context.Context, id string, limit int) ([]WebhookLog, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/webhooks/%s/logs", id), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if limit > 0 {
		q.Add("limit", strconv.Itoa(limit))
	}
	req.URL.RawQuery = q.Encode()

	var res []WebhookLog
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *WebhooksResource) Retry(ctx context.Context, logID string) error {
	req, err := r.c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/webhooks/logs/%s/retry", logID), map[string]interface{}{})
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}

func (r *WebhooksResource) Delete(ctx context.Context, id string) error {
	req, err := r.c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/webhooks/%s", id), nil)
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}
