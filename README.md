# Database Comparator [V0.0.1]

### Current Go version: 1.17

## Description:

This program has been created to compare two databases.
It shows the tables that are in common between both databases
and the tables which are only in one.
Then, it shows for every table in common between both databases the columns which
are shared between the tables of both databases, the columns which are only in one table, the ones which are present in
both tables but does not have the same data type and the ones which are present in both tables but does not have the
same data length.

## Database types currently supported:

- Postgresql (program argument `--type_[1|2] postgres`)
- Mysql (program argument `--type_[1|2] mysql`)

*Note: Two databases of two different types can be compared in the same command*

## Program arguments:
### Structure:
Program arguments can be provided in these four syntaxes:
- `-k value`
- `-k=value`
- `--key value`
- `--key=value`

(For the password flag, only the long form --password is available.)

### Arguments syntaxe:
Flags ending by **_1** will be parsed for the first database, these ending by **_2** for the second one.
The order does not matter.
#### Full syntax:
```shell
./DatabaseComparator-X.Y.Z --type_1 <db 1 type> --database_1 <db 1 name> --user_1 <user 1> [--password_1 <password 1>] --host_1 <host 1> [--port_1 <port 1>] --type_2 <db 2 type> --database_2 <db 2 name> --user_2 <user 2> [--password_2 <password 2>] --host_2 <host 2> [--port_2 <port 2>]
```
#### Short syntax:
```shell
./DatabaseComparator-X.Y.Z -t_1 <db 1 type> -db_1 <db 1 name> -u_1 <user 1> [--password_1 <password 1>] -h_1 <host 1> [-p_1 <port 1>] -t_2 <db 2 type> -db_2 <db 2 name> -u_2 <user 2> [--password_2 <password 2>] -h_2 <host 2> [-p_2 <port 2>]
```

### Argument details;
- `type`: The type of the database (refer to [this section](#database-types-currently-supported))
- `database`: The name of the database
- `user`: The username to connect to the database
- `password`: (Optionnal) The password to connect to the database
- `host`: The address of the database
- `port`: (Optionnal) The port the database listen on. The default value is the default port for the database type used (e.g.: 5432 for PgSql)
