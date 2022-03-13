SELECT table_name
FROM information_schema.tables
WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
  AND table_type = 'BASE TABLE'
ORDER BY table_name;