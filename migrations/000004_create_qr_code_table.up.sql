CREATE TABLE IF NOT EXISTS qr_code (
    id UUID PRIMARY KEY,
    owner_id UUID NOT NULL REFERENCES users(id),
    equipment_id UUID NOT NULL REFERENCES equipment (id),
    finance_id UUID NOT NULL REFERENCES finance (id),
    purchase_date DATE NOT NULL,
    residual_value NUMERIC(14, 2) NOT NULL,
    refund_amount NUMERIC(14, 2) NOT NULL
)