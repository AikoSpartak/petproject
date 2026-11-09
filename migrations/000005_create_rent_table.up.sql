CREATE TABLE IF NOT EXISTS rent (
    uuid UUID PRIMARY KEY,
    owner_uuid UUID NOT NULL,
    equipment_uuid VARCHAR[] DEFAULT '{}',
    orders_uuid UUID NOT NULL,
    branch_uuid UUID NOT NULL,
    created_at DATE NOT NULL,
    returns_date DATE NOT NULL
)
CREATE INDEX idx_owner_uuid ON rent (owner_uuid);
CREATE INDEX idx_equipment_uuid ON rent (equipment_uuid);
CREATE INDEX idx_branch_uuid ON rent (branch_uuid);
CREATE INDEX idx_created_at ON rent (created_at);
CREATE INDEX idx_returns_date ON rent (returns_date);



