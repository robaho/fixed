//go:build sql_scanner
// +build sql_scanner

package fixed

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSQL(t *testing.T) {
	// open a database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	// close the database at the end of the function
	defer db.Close()

	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS "test" ("id" INTEGER, "value" BLOB, PRIMARY KEY("id"))`); err != nil {
		t.Fatal(err)
	}
	v := NewF(1.2345)
	_, err = db.Exec("insert into test (value) values (?);", v)
	if err != nil {
		t.Fatal(err)
	}

	// get the customer record
	rows, err := db.Query("select value from test;")
	if err != nil {
		t.Fatal(err)
	}
	// close the rows at the end of the function
	defer rows.Close()
	var foundVal Fixed
	for rows.Next() {
		if err := rows.Scan(
			&foundVal,
		); err != nil {
			t.Fatal(err)
		}
		t.Log(foundVal)
		break
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if foundVal.String() != v.String() {
		t.Error("should be equal", foundVal.String(), v.String())
	}
}
