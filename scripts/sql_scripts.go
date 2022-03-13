package scripts

import _ "embed"

//go:embed pgsql/list_tables.sql
var PgsqlListTables string

//go:embed mysql/list_tables.sql
var MysqlListTables string
