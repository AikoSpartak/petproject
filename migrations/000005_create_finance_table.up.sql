CREATE TABLE IF NOT EXISTS finance (
    id UUID PRIMARY KEY,
    equipment_id UUID NOT NULL REFERENCES equipment (id),
    purchase_amount NUMERIC(14,2) NOT NULL,
    accrued_cost_of_use NUMERIC(14, 2) NOT NULL,
    refund_amount NUMERIC(14, 2) NOT NULL,
    withholding_amount NUMERIC(14 ,2) NOT NULL,
    logistics_fees NUMERIC(14, 2) NOT NULL,
)