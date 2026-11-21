CREATE TABLE IF NOT EXISTS rent (
    uuid UUID PRIMARY KEY,
    owner_uuid UUID NOT NULL,
    equipment_uuids VARCHAR[] DEFAULT '{}',
    branch_uuid UUID NOT NULL,
    employee_uuid UUID NOT NULL,
    invoices_uuid UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    rental_days INTEGER NOT NULL,
    returned_at TIMESTAMPTZ,
    is_returned BOOLEAN NOT NULL DEFAULT false,
    total_amount NUMERIC(10, 2) NOT NULL
);
CREATE INDEX idx_rent_owner_uuid ON rent (owner_uuid);
CREATE INDEX idx_rent_branch_uuid ON rent (branch_uuid);
CREATE INDEX idx_rent_created_at ON rent (created_at);
CREATE INDEX idx_rent_returned_at ON rent (returned_at);







