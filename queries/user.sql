DROP TABLE IF EXISTS users;

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  client_id TEXT UNIQUE NOT NULL,
  email TEXT UNIQUE NOT NULL,
  name TEXT NOT NULL,
  picture TEXT NOT NULL,
  created_at TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert Query
-- INSERT INTO users(client_id,email,name,picture,created_at) 
-- SELECT '109232021077514978643',
-- 'shdipaa73@gmail.com',
-- 'dipaa sh',
-- 'https://lh3.googleuserscontent.com/a/ACg8ocLVdradDUR6djzUMHZOdjFBSuy0xOxqqHTLj5o8JeNH=s96-c',
-- NOW()
-- WHERE NOT EXISTS (
-- SELECT client_id FROM users WHERE client_id = '109232021077514978643'
-- );
-- RETURNING id

-- INSERT user but get ID in either of cases
-- WITH s AS(
--   SELECT id FROM users WHERE client_id = '109232021077514978643'
-- ), 
-- ns AS(
--   INSERT INTO users(client_id,email,name,picture,created_at) 
--   SELECT '109232021077514978643',
--   'shdipaa73@gmail.com',
--   'dipaa sh',
--   'https://lh3.googleuserscontent.com/a/ACg8ocLVdradDUR6djzUMHZOdjFBSuy0xOxqqHTLj5o8JeNH=s96-c',
--   NOW()
--   WHERE NOT EXISTS (
--     SELECT 1 from s
--   )
--   RETURNING id
-- )
-- select id
-- from ns
-- union all
-- select id
-- from s
