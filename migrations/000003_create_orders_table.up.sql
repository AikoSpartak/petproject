-- Сделать индекс на эквипмент, бранч, селлер, креате_ат, 

CREATE TABLE IF NOT EXISTS orders (
    uuid UUID PRIMARY KEY,
    owner_uuid   UUID NOT NULL,
    equipment_uuid  VARCHAR[] DEFAULT [],
    branch_uuid UUID NOT NULL,
    seller_uuid UUID NOT NULL,
    created_at DATE NOT NULL,
    total_amount NUMERIC(10, 2) NOT NULL,
    paid BOOLEAN DEFAULT false
);

CREATE INDEX idx_equipment_uuid ON orders (equipment_uuid);
CREATE INDEX idx_branche_uuid ON orders (branche_uuid);
CREATE INDEX idx_seller_uuid ON orders (seller_uuid);
CREATE INDEX idx_created_at on orders (created_at);
