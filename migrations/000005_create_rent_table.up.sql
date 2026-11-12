CREATE TABLE IF NOT EXISTS rent (
    uuid UUID PRIMARY KEY,
    owner_uuid UUID NOT NULL,
    equipment_uuid VARCHAR[] DEFAULT [],
    branch_uuid UUID NOT NULL,
    created_at DATE NOT NULL,
    rental_days INTEGER NOT NULL CHECK (rental_days > 0),
    return_date DATE NOT NULL,
    returned BOOLEAN NOT NULL DEFAULT false,
    total_amount NUMERIC(10, 2) NOT NULL
)
CREATE INDEX idx_owner_uuid ON rent (owner_uuid);
CREATE INDEX idx_branch_uuid ON rent (branch_uuid);
CREATE INDEX idx_created_at ON rent (created_at);
CREATE INDEX idx_return_date ON rent (return_date);
CREATE INDEX idx_total_amount on rent (total_amount);



