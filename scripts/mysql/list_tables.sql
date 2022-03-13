SELECT TABLE_NAME as table_name
FROM information_schema.TABLES
WHERE Table_type = 'BASE TABLE'
ORDER BY TABLE_NAME;