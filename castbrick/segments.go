package castbrick

import (
	"context"
	"net/url"
	"strconv"
)

type SegmentsResource struct {
	c *Client
}

type ListSegmentsOptions struct {
	Page     int
	PageSize int
	Search   *string
}

func (r *SegmentsResource) List(ctx context.Context, opts ListSegmentsOptions) (*PagedResult[Segment], error) {
	q := url.Values{}
	if opts.Page > 0 {
		q.Set("pageNumber", strconv.Itoa(opts.Page))
	}
	if opts.PageSize > 0 {
		q.Set("pageSize", strconv.Itoa(opts.PageSize))
	}
	if opts.Search != nil && *opts.Search != "" {
		q.Set("search", *opts.Search)
	}

	var res PagedResult[Segment]
	if err := r.c.getQuery(ctx, "/audience/segments", q, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *SegmentsResource) Create(ctx context.Context, data CreateSegmentRequest) (string, error) {
	var res string
	if err := r.c.post(ctx, "/audience/segments", data, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *SegmentsResource) Update(ctx context.Context, id string, data UpdateSegmentRequest) error {
	path, err := pathf("/audience/segments/%s", id)
	if err != nil {
		return err
	}
	return r.c.put(ctx, path, data, nil)
}

func (r *SegmentsResource) Delete(ctx context.Context, id string) error {
	path, err := pathf("/audience/segments/%s", id)
	if err != nil {
		return err
	}
	return r.c.delete(ctx, path)
}
