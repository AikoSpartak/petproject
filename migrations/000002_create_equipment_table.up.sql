CREATE TYPE status_equipment AS ENUM (
    'На складе',
    'В прокате',
    'Продан с обратым выкупом',
    'Возвращен',
    'На ремонте',
    'Списан',
    'Перемещается между филиалами'
)



CREATE TABLE IF NOT EXISTS equipment (
    id                UUID PRIMARY KEY,
    ownew_id          UUID NOT NULL REFERENCES users(id),
    status            status_equipment NOT NULL DEFAULT 'На складе',
    home_branch_id    UUID NOT NULL REFERENCES branches(id),  
    current_branch_id UUID NOT NULL REFERENCES branches(id),  
    category          VARCHAR(256) NOT NULL,
    brand             VARCHAR(256) NOT NULL,
    model             VARCHAR(256) NOT NULL,
    size_equipment    NUMERIC(4,1) NOT NULL,
    length_equipment  NUMERIC(4,1) NOT NULL,
    serial_number     VARCHAR(256) UNIQUE NOT NULL
);