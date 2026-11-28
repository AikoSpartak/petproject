CREATE TABLE IF NOT EXISTS returs (
    uuid UUID PRIMARY KEY,
    owner_uuid UUID NOT NULL,
    equipment_uuid UUID NOT NULL,
    branch_uuid UUID NOT NULL,
    is_returned BOOLEAN NOT NULL DEFAULT false
)