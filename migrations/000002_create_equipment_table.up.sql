CREATE TABLE IF NOT EXISTS equipment (
    uuid              UUID PRIMARY KEY,
    status            VARCHAR(256) DEFAULT 'На складе',
    owner_uuid        UUID,
    branch_uuid       UUID,
    description       TEXT, 
    category          VARCHAR(256) NOT NULL,
    brand             VARCHAR(256) NOT NULL,
    model             VARCHAR(256) NOT NULL,
    specification     JSON DEFAULT '{}',
    serial_number     VARCHAR(256) UNIQUE NOT NULL,
    price             NUMERIC(10,2) NOT NULL DEFAULT 0,
    rend_rate         JSON DEFAULT '{}'
);

CREATE INDEX idx_status ON equipment (status);
CREATE INDEX idx_owner_uuid ON equipment (owner_uuid);
CREATE INDEX idx_branch_uuid ON equipment (branch_uuid);
CREATE INDEX idx_category ON equipment (category);
CREATE INDEX idx_model ON equipment (model);
CREATE INDEX idx_brand ON equipment (brand);
CREATE INDEX idx_serial_number ON equipment (serial_number); 