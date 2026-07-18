package castbrick

import (
	"context"
	"fmt"
	"net/http"
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
	req, err := r.c.newRequest(ctx, http.MethodGet, "/audience/segments", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Page > 0 {
		q.Add("pageNumber", strconv.Itoa(opts.Page))
	}
	if opts.PageSize > 0 {
		q.Add("pageSize", strconv.Itoa(opts.PageSize))
	}
	if opts.Search != nil && *opts.Search != "" {
		q.Add("search", *opts.Search)
	}
	req.URL.RawQuery = q.Encode()

	var res PagedResult[Segment]
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *SegmentsResource) Create(ctx context.Context, data CreateSegmentRequest) (string, error) {
	req, err := r.c.newRequest(ctx, http.MethodPost, "/audience/segments", data)
	if err != nil {
		return "", err
	}

	var res string
	if err := r.c.do(req, &res); err != nil {
		return "", err
	}
	return res, nil
}

func (r *SegmentsResource) Update(ctx context.Context, id string, data UpdateSegmentRequest) error {
	req, err := r.c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/audience/segments/%s", id), data)
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}

func (r *SegmentsResource) Delete(ctx context.Context, id string) error {
	req, err := r.c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/audience/segments/%s", id), nil)
	if err != nil {
		return err
	}

	return r.c.do(req, nil)
}
