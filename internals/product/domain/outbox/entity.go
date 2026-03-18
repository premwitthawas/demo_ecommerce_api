package product

import "time"

type ProductOutboxMessageStatus string

const (
	PENDING   ProductOutboxMessageStatus = "pending"
	PUBLISHED ProductOutboxMessageStatus = "published"
	RETRYING  ProductOutboxMessageStatus = "retrying"
	DLQ       ProductOutboxMessageStatus = "dql"
)

type ProductOutboxMessage struct {
	ID          string     `json:"id"`
	EventType   string     `json:"event_type"`
	AggrID      string     `json:"aggr_id"`
	AggrVersion int32      `json:"aggr_version"`
	Status      string     `json:"status"`
	Payload     []byte     `json:"payload"`
	Metadata    []byte     `json:"metadata"`
	RetryCount  int32      `json:"retry_count"`
	NextRetryAt time.Time  `json:"next_retry_at"`
	ErrText     *string    `json:"err_text"`
	ConsumedAt  *time.Time `json:"consumed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Version     int32      `json:"version"`
}

func NewProductOutbox(payload *ProductOutboxMessage) (*ProductOutboxMessage, error) {
	if payload.ID == "" {
		return nil, ErrProductOutboxIDEmpty
	}
	if payload.AggrID == "" {
		return nil, ErrProductOutboxAggrIDEmpty
	}
	if payload.EventType == "" {
		return nil, ErrProductOutboxEventTypeEmpty
	}
	if payload.Status == "" {
		return nil, ErrProductOutboxStatustypeEmpty
	}
	if len(payload.Payload) == 0 {
		return nil, ErrProductOutboxPayloadEmpty
	}
	if len(payload.Metadata) == 0 {
		return nil, ErrProductOutboxMetadataEmpty
	}
	now := time.Now()
	return &ProductOutboxMessage{
		ID:          payload.ID,
		EventType:   payload.EventType,
		AggrID:      payload.AggrID,
		AggrVersion: payload.AggrVersion,
		Status:      string(PENDING),
		Payload:     payload.Payload,
		Metadata:    payload.Metadata,
		RetryCount:  payload.RetryCount,
		NextRetryAt: now,
		CreatedAt:   now,
		UpdatedAt:   now,
		Version:     1,
		ErrText:     nil,
		ConsumedAt:  nil,
	}, nil
}
