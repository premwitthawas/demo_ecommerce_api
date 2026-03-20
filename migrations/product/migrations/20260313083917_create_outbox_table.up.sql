CREATE TABLE IF NOT EXISTS outbox_messages (
    id varchar(255) primary key,
    event_type varchar(255) not null,
    aggr_id varchar(255) not null,
    aggr_version int not null,
    status varchar(50) not null,
    payload jsonb,
    metadata jsonb,
    retry_count int not null default 0,
    next_retry_at timestamptz not null,
    err_text text,
    consumed_at timestamptz,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),
    version int not null default 1
);

CREATE INDEX idx_product_outbox_message_polling on outbox_messages(next_retry_at)
WHERE status = 'pending';
