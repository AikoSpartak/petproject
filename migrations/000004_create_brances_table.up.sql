CREATE TABLE IF NOT EXISTS branches (
    uuid UUID PRIMARY KEY,
    address TEXT NOT NULL,
    description TEXT NOT NULL
)