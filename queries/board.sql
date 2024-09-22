DROP TABLE IF EXISTS board;

CREATE TABLE board(
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(255),
    slug VARCHAR(255) UNIQUE NOT NULL,
    data JSONB,
    owner UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    space_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (id,space_id)
)
