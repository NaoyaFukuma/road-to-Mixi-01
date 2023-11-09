package testhelpers

import (
	"database/sql"
	"fmt"
	"testing"
)

// Dump DB is a test helper function to dump the contents of a database table.
func DumpDB(t *testing.T, db *sql.DB, tableName string) {
	t.Helper()
	t.Logf("Dumping table %s", tableName)
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		t.Fatalf("failed to query table %s: %v", tableName, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		t.Fatalf("failed to get columns for table %s: %v", tableName, err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			t.Fatalf("failed to scan row: %v", err)
		}

		for i, col := range columns {
			val := values[i]

			b, ok := val.([]byte)
			if ok {
				t.Logf("%s: %s\n", col, string(b))
			} else {
				t.Logf("%s: %v\n", col, val)
			}
		}
	}
}
