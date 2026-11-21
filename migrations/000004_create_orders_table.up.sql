CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY,
    owner_uuid   UUID NOT NULL,
    equipment_uuid  VARCHAR[] DEFAULT '{}',
    branch_uuid UUID NOT NULL,
    seller_uuid UUID NOT NULL,
    invoices_uuid UUID NOT NULL,
    created_at DATE NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    is_paid BOOLEAN DEFAULT false
);

CREATE INDEX idx_equipment_uuid ON orders (equipment_uuid);
CREATE INDEX idx_branche_uuid ON orders (branch_uuid);
CREATE INDEX idx_seller_uuid ON orders (seller_uuid);
CREATE INDEX idx_created_at ON orders (created_at);
