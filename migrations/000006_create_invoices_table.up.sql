CREATE TABLE IF NOT EXISTS invoices (
    uuid           UUID PRIMARY KEY,
    customer_uuid  UUID NOT NULL,
    employee_uuid  UUID NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    paid_at        TIMESTAMPTZ,
    amount         NUMERIC(10,2) NOT NULL CHECK (amount > 0),
    is_rejected    BOOLEAN NOT NULL DEFAULT false,
    is_paid        BOOLEAN NOT NULL DEFAULT false
);

CREATE INDEX idx_invoices_customer_uuid ON invoices (customer_uuid);
CREATE INDEX idx_invoices_employee_uuid ON invoices (employee_uuid);
CREATE INDEX idx_invoices_created_at ON invoices (created_at);
