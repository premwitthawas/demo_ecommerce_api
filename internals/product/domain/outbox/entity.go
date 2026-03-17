package product

import "time"

type OutboxMessage struct {
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
	Version     *int32     `json:"version"`
}
