package castbrick

import "time"

// PagedResult is a generic paginated response.
type PagedResult[T any] struct {
	Items          []T  `json:"items"`
	TotalCount     int  `json:"totalCount"`
	PageNumber     int  `json:"pageNumber"`
	TotalPages     int  `json:"totalPages"`
	HasNextPage    bool `json:"hasNextPage"`
	HasPreviousPage bool `json:"hasPreviousPage"`
}

// ── SMS ──────────────────────────────────────────────────────────────────────

type SendSmsResponse struct {
	MessageID      string `json:"messageId"`
	Status         string `json:"status"`
	RecipientCount int    `json:"recipientCount"`
	Error          string `json:"error,omitempty"`
	Timestamp      string `json:"timestamp"`
}

type SmsMessage struct {
	ID            string     `json:"id"`
	ContactName   *string    `json:"contactName,omitempty"`
	RecipientPhone string    `json:"recipientPhone"`
	Message       string     `json:"message"`
	CampaignName  *string    `json:"campaignName,omitempty"`
	CampaignID    *string    `json:"campaignId,omitempty"`
	SenderID      *string    `json:"senderId,omitempty"`
	Status        string     `json:"status"`
	ErrorMessage  *string    `json:"errorMessage,omitempty"`
	RetryCount    int        `json:"retryCount"`
	ScheduledAt   *time.Time `json:"scheduledAt,omitempty"`
	SentAt        *time.Time `json:"sentAt,omitempty"`
	DeliveredAt   *time.Time `json:"deliveredAt,omitempty"`
}

// ── Contacts ─────────────────────────────────────────────────────────────────

type Contact struct {
	ID          string    `json:"id"`
	Name        string    `json:"name,omitempty"`
	PhoneNumber string    `json:"phoneNumber,omitempty"`
	Email       string    `json:"email,omitempty"`
	TenantID    string    `json:"tenantId"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ContactList struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	TenantID     string    `json:"tenantId"`
	ContactCount int       `json:"contactCount"`
	CreatedAt    time.Time `json:"createdAt"`
}

// ── Broadcasts ───────────────────────────────────────────────────────────────

type Broadcast struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Status        string     `json:"status"`
	Message       string     `json:"message"`
	SenderID      string     `json:"senderId,omitempty"`
	ContactListID string     `json:"contactListId,omitempty"`
	ScheduledAt   *time.Time `json:"scheduledAt,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
}

// ── Templates ────────────────────────────────────────────────────────────────

type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	Subject   *string   `json:"subject,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateTemplateRequest struct {
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Subject *string `json:"subject,omitempty"`
}

type UpdateTemplateRequest struct {
	Name    string  `json:"name"`
	Content string  `json:"content"`
	Subject *string `json:"subject,omitempty"`
}

// ── Webhooks ─────────────────────────────────────────────────────────────────

type Webhook struct {
	ID        string    `json:"id"`
	Endpoint  string    `json:"endpoint"`
	EventType string    `json:"eventType"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateWebhookRequest struct {
	Endpoint  string `json:"endpoint"`
	EventType string `json:"eventType"`
}

type WebhookLog struct {
	ID          string    `json:"id"`
	WebhookID   string    `json:"webhookId"`
	EventType   string    `json:"eventType"`
	StatusCode  int       `json:"statusCode"`
	Attempt     int       `json:"attempt"`
	DeliveredAt time.Time `json:"deliveredAt"`
}

// ── Segments ─────────────────────────────────────────────────────────────────

type Segment struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   *string   `json:"description,omitempty"`
	RulesOperator *string   `json:"rulesOperator,omitempty"`
	ContactCount  int       `json:"contactCount"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateSegmentRequest struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	RulesOperator *string `json:"rulesOperator,omitempty"`
	Rules         *string `json:"rules,omitempty"`
}

type UpdateSegmentRequest struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	RulesOperator *string `json:"rulesOperator,omitempty"`
	Rules         *string `json:"rules,omitempty"`
}

// ── Billing ──────────────────────────────────────────────────────────────────

type BillingBalance struct {
	Balance             float64    `json:"balance"`
	LowBalanceThreshold float64    `json:"lowBalanceThreshold"`
	AlertEmail          string     `json:"alertEmail"`
	LastAlertSentAt     *time.Time `json:"lastAlertSentAt,omitempty"`
}

type CreditPack struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Credits     int     `json:"credits"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	PricePerSms float64 `json:"pricePerSms"`
}

type InitiatePaymentRequest struct {
	PackageID  string `json:"packageId"`
	Provider   string `json:"provider"`
	Method     string `json:"method"`
	Currency   string `json:"currency"`
	ReturnURL  string `json:"returnUrl"`
	CancelURL  string `json:"cancelUrl"`
	SuccessURL string `json:"successUrl"`
	FailureURL string `json:"failureUrl"`
}

type InitiatePaymentResponse struct {
	PaymentURL  string `json:"paymentUrl"`
	ProviderRef string `json:"providerRef"`
	PaymentID   string `json:"paymentId"`
}

type PaymentInfo struct {
	ID             string    `json:"id"`
	Amount         float64   `json:"amount"`
	Currency       string    `json:"currency"`
	Status         string    `json:"status"`
	ProviderRef    string    `json:"providerRef"`
	CreditsGranted int       `json:"creditsGranted"`
	CreatedAt      time.Time `json:"createdAt"`
}

type CreditLedgerEntry struct {
	ID          string    `json:"id"`
	Amount      int       `json:"amount"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PaymentListResponse struct {
	Items      []PaymentInfo `json:"items"`
	NextCursor *string       `json:"nextCursor,omitempty"`
}

type TransactionListResponse struct {
	Items      []CreditLedgerEntry `json:"items"`
	NextCursor *string             `json:"nextCursor,omitempty"`
}
