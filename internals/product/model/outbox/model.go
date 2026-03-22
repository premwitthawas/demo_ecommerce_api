package product_outbox

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

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

func calulateJitter() time.Time {
	ms := rand.Intn(500)
	return time.Now().Add(time.Duration(ms) * time.Millisecond)
}

func (o *ProductOutboxMessage) UpdateOutboxPublished() {
	now := time.Now()
	// o.Version++
	o.UpdatedAt = now
	o.ConsumedAt = &now
	o.Status = string(PUBLISHED)
}

func (o *ProductOutboxMessage) UpdateOutboxRetrying(err error) {
	now := time.Now()
	// o.Version++
	o.RetryCount++
	backoff := time.Duration(math.Pow(2, float64(o.RetryCount))) * time.Second
	o.NextRetryAt = calulateJitter().Add(backoff)
	if err != nil {
		errStr := err.Error()
		o.ErrText = &errStr
	} else {
		defaultErr := "unknown error during retry"
		o.ErrText = &defaultErr
	}
	o.UpdatedAt = now
	o.Status = string(RETRYING)
}

func (o *ProductOutboxMessage) UpdateOutboxDLQ(err error) {
	now := time.Now()
	// o.Version++
	o.RetryCount++
	errStr := fmt.Sprintf("Final attempt failed: %v", err)
	o.ErrText = &errStr
	o.UpdatedAt = now
	o.Status = string(DLQ)
}
