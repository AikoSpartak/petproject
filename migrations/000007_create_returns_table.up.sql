
CREATE TABLE IF NOT EXISTS returns_equipment (
    uuid UUID PRIMARY KEY,
    owner_uuid UUID NOT NULL,
    equipment_uuid VARCHAR[] DEFAULT '{}',
    state_equipment VARCHAR(256) DEFAULT 'отличное',
    created_at DATE NOT NULL,
    sum_returns NUMERIC(10, 2)
)

CREATE INDEX idx_owner_uuid ON returns_equipment (owner_uui);
CREATE INDEX idx_state_equipment ON returns_equipment (state_equipment);
