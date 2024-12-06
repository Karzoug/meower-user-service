CREATE TABLE outbox (
    id SERIAL,
    change_type VARCHAR(100) NOT NULL,
    user_id public.xid NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reserved_to TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id)
);