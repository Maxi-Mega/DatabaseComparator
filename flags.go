package main

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var flagTagRegexp = regexp.MustCompile(`flag:"(?P<flag>[-\w,]+)"`)
var flagEqualRawRegexp = `%s=(?P<value>\w+)`

// parseConnType looks for the type flag in the program arguments and returns the corresponding connectionType,
// or crashes if no type is provided or if it is invalid
func parseConnType(args *[]string, connId string) (connType connectionType) {
	dbType, found := findFlag(args, "type", connId, "-t", "--type")
	if !found {
		log.Fatalf("No database type found: missing flag -t_%[1]s/--type_%[1]s", connId)
	}
	switch strings.ToLower(dbType) {
	case "postgres":
		connType = &pgsqlConn{}
	case "mysql":
		connType = &mysqlConn{}
	default:
		log.Fatalf("Unknown database type: %q\n", dbType)
	}
	if err := connType.parseArgs(args, connId); err != nil {
		log.Fatalln("Failed to parse args for database type", dbType+":", err)
	}
	return
}

// fillStructFrom tries to fill every field of the given connectionType with the values contained in the given args
func fillStructFrom(args *[]string, conn connectionType, connId string) error {
	val := reflect.ValueOf(conn)
	if val.Kind() != reflect.Ptr {
		return errors.New("given connection is not a pointer")
	}
	s := val.Elem()
	if s.Kind() != reflect.Struct {
		return errors.New("given connection is not a pointer")
	}
	fieldsNbr := s.Type().NumField()
	for i := 0; i < fieldsNbr; i++ {
		f := s.Field(i)
		if f.IsValid() {
			field := reflect.TypeOf(conn).Elem().Field(i)
			if !field.IsExported() {
				continue
			}
			// extracting the list of flag aliases from the field's tag
			flagNames, found := findGroup(flagTagRegexp, string(field.Tag), "flag")
			if !found {
				continue
			}
			flagValue, found := parseFlag(strings.Split(flagNames, ","), args, f.Kind(), connId)
			if !found {
				continue
			}
			if f.CanSet() {
				f.Set(reflect.ValueOf(flagValue))
			}
		}
	}
	return nil
}

// findGroup looks for the given groupName in the matches of the given string with the given regexp
func findGroup(re *regexp.Regexp, str, groupName string) (value string, found bool) {
	matches := re.FindStringSubmatch(str)
	if matches != nil {
		for gi, gn := range re.SubexpNames() {
			if gi != 0 && gn == groupName {
				return matches[gi], true
			}
		}
	}
	return "", false
}

// parseFlag looks for the flag with the given flagNames in the given args and tries to convert it into the expectedType
func parseFlag(flagNames []string, args *[]string, expectedType reflect.Kind, connId string) (value interface{}, found bool) {
	val, found := findFlag(args, "value", connId, flagNames...)
	if !found {
		return
	}
	switch expectedType {
	case reflect.String:
		return val, true
	case reflect.Int:
		i, err := strconv.Atoi(val)
		if err != nil {
			return nil, false
		}
		return i, true
	case reflect.Uint:
		i, err := strconv.ParseUint(val, 10, 32)
		if err != nil {
			return nil, false
		}
		return uint(i), true
	default:
		return nil, false
	}
}

// findFlag looks for the flag with the given flagNames in the given args
func findFlag(args *[]string, reGroupName, connId string, flagNames ...string) (value string, found bool) {
	valueIsNextArg := false
	for i, arg := range *args {
		if valueIsNextArg {
			*args = append((*args)[:i-1], (*args)[i+1:]...)
			return arg, true
		}

		for _, flagName := range flagNames {
			flagName += "_" + connId
			re := regexp.MustCompile(fmt.Sprintf(flagEqualRawRegexp, flagName))
			if value, found = findGroup(re, arg, reGroupName); found {
				*args = append((*args)[:i], (*args)[i+1:]...)
				return value, true
			}
			if arg == flagName {
				valueIsNextArg = true
				continue
			}
		}
	}
	return "", false
}
