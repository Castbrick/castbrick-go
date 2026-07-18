package castbrick

import (
	"context"
	"net/http"
	"strconv"
)

type BillingResource struct {
	c *Client
}

func (r *BillingResource) GetBalance(ctx context.Context) (*BillingBalance, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/billing/balance", nil)
	if err != nil {
		return nil, err
	}

	var res BillingBalance
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *BillingResource) ListPackages(ctx context.Context) ([]CreditPack, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/billing/packages", nil)
	if err != nil {
		return nil, err
	}

	var res []CreditPack
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r *BillingResource) InitiatePayment(ctx context.Context, data InitiatePaymentRequest) (*InitiatePaymentResponse, error) {
	req, err := r.c.newRequest(ctx, http.MethodPost, "/billing/payments/initiate", data)
	if err != nil {
		return nil, err
	}

	var res InitiatePaymentResponse
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type ListPaymentsOptions struct {
	Cursor *string
	Limit  int
	Status *string
	From   *string
	To     *string
}

func (r *BillingResource) ListPayments(ctx context.Context, opts ListPaymentsOptions) (*PaymentListResponse, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/billing/payments", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Cursor != nil && *opts.Cursor != "" {
		q.Add("cursor", *opts.Cursor)
	}
	if opts.Limit > 0 {
		q.Add("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Status != nil && *opts.Status != "" {
		q.Add("status", *opts.Status)
	}
	if opts.From != nil && *opts.From != "" {
		q.Add("from", *opts.From)
	}
	if opts.To != nil && *opts.To != "" {
		q.Add("to", *opts.To)
	}
	req.URL.RawQuery = q.Encode()

	var res PaymentListResponse
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type ListTransactionsOptions struct {
	Cursor *string
}

func (r *BillingResource) ListTransactions(ctx context.Context, opts ListTransactionsOptions) (*TransactionListResponse, error) {
	req, err := r.c.newRequest(ctx, http.MethodGet, "/billing/transactions", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if opts.Cursor != nil && *opts.Cursor != "" {
		q.Add("cursor", *opts.Cursor)
	}
	req.URL.RawQuery = q.Encode()

	var res TransactionListResponse
	if err := r.c.do(req, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
