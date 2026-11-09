CREATE TABLE IF NOT EXISTS invoices (
    uuid UUID PRIMARY KEY,
    users_uuid UUID NOT NULL,
    created_at DATE NOT NULL, 
    equipment_uuid VARCHAR[] DEFAULT '{}',
    total_sum NUMERIC(10, 2) NOT NULL,
    change NUMERIC(10, 2) NOT NULL
)

CREATE INDEX idx_created_at ON invoices (created_at);
CREATE INDEX idx_total_sum ON invoices (total_sum);
