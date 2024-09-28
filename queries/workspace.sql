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

-- Workspace create :
-- Insert into workspace table -> then insert into members table.

-- WITH spaceins AS (
--     INSERT INTO workspace (name,overview,slug,owner,created_at)
--     SELECT 'Icebox','Workspace for players to understand Icebox map','icebox','306b199b-5429-4826-b568-a55ebb26deba',NOW()
--     RETURNING id,name,overview,slug,owner,created_at,updated_at
-- ),
-- memins AS (
--     INSERT INTO members (user_id,space_id)
--     SELECT '9b1f5ebd-b3c4-42e1-a1ef-f2c20d6a74a7',id
-- )
-- SELECT id,name,overview,slug,owner,created_at,updated_at FROM spaceins;


-- WITH spaceins AS (
--     INSERT INTO workspace (name,overview,slug,owner,created_at)
--     SELECT 'Icebox','Workspace for players to understand Icebox map','icebox','306b199b-5429-4826-b568-a55ebb26deba',NOW()
--     RETURNING *
-- ),
-- memins AS (
--     INSERT INTO member (user_id,space_id)
--     VALUES 
--     ('9b1f5ebd-b3c4-42e1-a1ef-f2c20d6a74a7',(SELECT id from spaceins)),
--     ('306b199b-5429-4826-b568-a55ebb26deba',(SELECT id from spaceins))
-- )
-- SELECT id,name,overview,slug,owner,created_at,updated_at FROM spaceins;



-- SELECT 
--     w.*,
--     o.* AS ownerdetails,
--     json_agg(u) AS members
-- FROM 
--     workspace w
-- LEFT JOIN 
--     member m ON w.id = m.space_id
-- LEFT JOIN 
--     users u ON m.user_id = u.id
-- JOIN 
--     users o ON w.owner = o.id
-- GROUP BY 
--     w.id, o.id;
