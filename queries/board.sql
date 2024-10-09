DROP TABLE IF EXISTS board;

CREATE TABLE board(
    id UUID DEFAULT gen_random_uuid(),
    name VARCHAR(255) UNIQUE NOT NULL,
    type VARCHAR(255),
    slug VARCHAR(255) UNIQUE NOT NULL,
    data JSONB,
    owner UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    space_id UUID NOT NULL REFERENCES workspace(id) ON DELETE CASCADE,
    PRIMARY KEY (id,space_id)
)

-- Board create
-- INSERT INTO board (name,type,slug,owner,space_id)
-- SELECT 'Icebox kanban','kanban','icebox-kanban','306b199b-5429-4826-b568-a55ebb26deba','23d372a3-cc3a-4637-bbf9-83dcb0437934'
-- RETURNING *

