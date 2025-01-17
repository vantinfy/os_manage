package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"modernc.org/sqlite"
	"regexp"
	"strings"

	_ "modernc.org/sqlite"
)

func RegisterRegexp() {
	err := sqlite.RegisterDeterministicScalarFunction(
		"regexp",
		2,
		func(ctx *sqlite.FunctionContext, args []driver.Value) (driver.Value, error) {
			s1 := args[0].(string)
			s2 := args[1].(string)

			matched, err := regexp.MatchString(s1, s2)
			if err != nil {
				return nil, fmt.Errorf("bad regular expression: %v", err)
			}
			return matched, nil
		},
	)
	if err != nil {
		if strings.Contains(err.Error(), `a function named "regexp" is already registered`) {
			return
		}
		panic(fmt.Errorf("sqlite registration failed: %v", err))
	}
}

func GetDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	return db, err
}
