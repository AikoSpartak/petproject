CREATE TABLE IF NOT EXISTS invoices (
    uuid UUID PRIMARY KEY,
    user_uuid UUID NOT NULL,
    amount_to_pay NUMERIC(10, 2) NOT NULL CHECK (amount_to_pay > 0),
    paid BOOLEAN NOT NULL DEFAULT false,
)

CREATE INDEX idx_user_uuid ON invoices (user_uuid);
