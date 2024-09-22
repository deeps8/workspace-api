DROP TABLE IF EXISTS workspace;

CREATE TABLE workspace (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    overview TEXT NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    owner UUID NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_owner
        FOREIGN KEY(owner)
            REFERENCES users(id)
            ON DELETE CASCADE
)