package castbrick

import (
	"context"
	"net/url"
	"strconv"
)

type BillingResource struct {
	c *Client
}

func (r *BillingResource) GetBalance(ctx context.Context) (*BillingBalance, error) {
	var res BillingBalance
	if err := r.c.get(ctx, "/billing/balance", &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *BillingResource) ListPackages(ctx context.Context) ([]CreditPack, error) {
	var res []CreditPack
	if err := r.c.get(ctx, "/billing/packages", &res); err != nil {
		return nil, err
	}
	return res, nil
}

// InitiatePayment starts a payment for a credit package.
//
// The redirect URLs (ReturnURL, CancelURL, SuccessURL, FailureURL) are sent to
// the payment provider as-is. Never fill them from end-user input without
// checking them against an allowlist, or they become an open redirect.
func (r *BillingResource) InitiatePayment(ctx context.Context, data InitiatePaymentRequest) (*InitiatePaymentResponse, error) {
	var res InitiatePaymentResponse
	if err := r.c.post(ctx, "/billing/payments/initiate", data, &res); err != nil {
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
	q := url.Values{}
	if opts.Cursor != nil && *opts.Cursor != "" {
		q.Set("cursor", *opts.Cursor)
	}
	if opts.Limit > 0 {
		q.Set("limit", strconv.Itoa(opts.Limit))
	}
	if opts.Status != nil && *opts.Status != "" {
		q.Set("status", *opts.Status)
	}
	if opts.From != nil && *opts.From != "" {
		q.Set("from", *opts.From)
	}
	if opts.To != nil && *opts.To != "" {
		q.Set("to", *opts.To)
	}

	var res PaymentListResponse
	if err := r.c.getQuery(ctx, "/billing/payments", q, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

type ListTransactionsOptions struct {
	Cursor *string
}

func (r *BillingResource) ListTransactions(ctx context.Context, opts ListTransactionsOptions) (*TransactionListResponse, error) {
	q := url.Values{}
	if opts.Cursor != nil && *opts.Cursor != "" {
		q.Set("cursor", *opts.Cursor)
	}

	var res TransactionListResponse
	if err := r.c.getQuery(ctx, "/billing/transactions", q, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
