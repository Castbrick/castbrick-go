package castbrick

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type TemplatesResource struct {
	c *Client
}

func (r *TemplatesResource) List(ctx context.Context, page, pageSize int) (*PagedResult[Template], error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/templates", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if page > 0 {
		q.Add("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Add("pageSize", strconv.Itoa(pageSize))
	}
	req.URL.RawQuery = q.Encode()

	var res PagedResult[Template]
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *TemplatesResource) Create(ctx context.Context, data CreateTemplateRequest) (string, error) {
	req, err := r.c.newRequest(ctx, http.MethodPost, "/templates", data)
	if err != nil {
		return "", err
	}

	var res string
	if err := r.c.do(req, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *TemplatesResource) Update(ctx context.Context, id string, data UpdateTemplateRequest) (string, error) {
	req, err := r.c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/templates/%s", id), data)
	if err != nil {
		return "", err
	}

	var res string
	if err := r.c.do(req, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *TemplatesResource) Delete(ctx context.Context, id string) error {
	req, err := r.c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/templates/%s", id), nil)
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}
