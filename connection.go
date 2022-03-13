package main

import (
	"DatabaseComparator/scripts" // TODO: github link
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// connectionType is a wrapper for a sql driver
type connectionType interface {
	// parseArgs tries to fill the fields with the values contained in the given args
	parseArgs(args *[]string, connId string) error
	// getDriver returns the name of the connection driver
	getDriver() string
	// getConnString builds a valid connection string from the loaded parameters
	getConnString() string
	// getListTablesScript returns a sql script that will list all the tables contained in the database
	getListTablesScript() string
	// getListColumnsScript returns a sql script that will list all the columns and their properties of the given table
	getListColumnsScript(table string) string
}

type pgsqlConn struct {
	DbName   string `flag:"-db,--database"`
	User     string `flag:"-u,--user"`
	Password string `flag:"--password"`
	Host     string `flag:"-h,--host"`
	Port     uint   `flag:"-p,--port"`
}

// parseArgs tries to fill the fields with the values contained in the given args
func (conn *pgsqlConn) parseArgs(args *[]string, connId string) error {
	err := fillStructFrom(args, conn, connId) // give a map with default values ?
	if err != nil {
		return err
	}
	if conn.DbName == "" {
		return errors.New("no database name provided")
	}
	if conn.Host == "" {
		return errors.New("no host provided")
	}
	if conn.User == "" {
		return errors.New("no user provided")
	}
	if conn.Port == 0 {
		conn.Port = 5432
	}
	return nil
}

// getDriver returns the name of the connection driver
func (conn pgsqlConn) getDriver() string {
	return "postgres"
}

// getConnString builds a valid connection string from the loaded parameters
func (conn pgsqlConn) getConnString() string {
	if conn.Password != "" {
		return fmt.Sprintf("%s://%s:%s@%s:%d/%s", conn.getDriver(), conn.User, conn.Password, conn.Host, conn.Port, conn.DbName)
	}
	return fmt.Sprintf("%s://%s@%s:%d/%s", conn.getDriver(), conn.User, conn.Host, conn.Port, conn.DbName)
}

// getListTablesScript returns a sql script that will list all the tables contained in the database
func (conn pgsqlConn) getListTablesScript() string {
	return scripts.PgsqlListTables
}

// getListColumnsScript returns a sql script that will list all the columns and their properties of the given table
func (conn pgsqlConn) getListColumnsScript(table string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table)
}

type mysqlConn struct {
	DbName   string `flag:"-db,--database"`
	User     string `flag:"-u,--user"`
	Password string `flag:"--password"`
	Host     string `flag:"-h,--host"`
	Port     uint   `flag:"-p,--port"`
}

// parseArgs tries to fill the fields with the values contained in the given args
func (conn *mysqlConn) parseArgs(args *[]string, connId string) error {
	err := fillStructFrom(args, conn, connId) // give a map with default values ?
	if err != nil {
		return err
	}
	if conn.DbName == "" {
		return errors.New("no database name provided")
	}
	if conn.Host == "" {
		return errors.New("no host provided")
	}
	if conn.User == "" {
		return errors.New("no user provided")
	}
	if conn.Port == 0 {
		conn.Port = 3306
	}
	return nil
}

// getDriver returns the name of the connection driver
func (conn mysqlConn) getDriver() string {
	return "mysql"
}

// getConnString builds a valid connection string from the loaded parameters
func (conn mysqlConn) getConnString() string {
	if conn.Password != "" {
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conn.User, conn.Password, conn.Host, conn.Port, conn.DbName)
	}
	return fmt.Sprintf("%s@tcp(%s:%d)/%s", conn.User, conn.Host, conn.Port, conn.DbName)
}

// getListTablesScript returns a sql script that will list all the tables contained in the database
func (conn mysqlConn) getListTablesScript() string {
	return scripts.MysqlListTables
}

// getListColumnsScript returns a sql script that will list all the columns and their properties of the given table
func (conn mysqlConn) getListColumnsScript(table string) string {
	return fmt.Sprintf("SELECT * FROM %s LIMIT 1;", table)
}
