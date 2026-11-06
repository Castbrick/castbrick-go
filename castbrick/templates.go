package castbrick

import (
	"context"
	"net/url"
	"strconv"
)

type TemplatesResource struct {
	c *Client
}

func (r *TemplatesResource) List(ctx context.Context, page, pageSize int) (*PagedResult[Template], error) {
	q := url.Values{}
	if page > 0 {
		q.Set("page", strconv.Itoa(page))
	}
	if pageSize > 0 {
		q.Set("pageSize", strconv.Itoa(pageSize))
	}

	var res PagedResult[Template]
	if err := r.c.getQuery(ctx, "/templates", q, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *TemplatesResource) Create(ctx context.Context, data CreateTemplateRequest) (string, error) {
	var res string
	if err := r.c.post(ctx, "/templates", data, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *TemplatesResource) Update(ctx context.Context, id string, data UpdateTemplateRequest) (string, error) {
	path, err := pathf("/templates/%s", id)
	if err != nil {
		return "", err
	}
	var res string
	if err := r.c.put(ctx, path, data, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *TemplatesResource) Delete(ctx context.Context, id string) error {
	path, err := pathf("/templates/%s", id)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
