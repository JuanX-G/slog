package database

import (
	"context"
	"fmt"
)

func (d *DB) InsertInto(ctx context.Context, tableName string, cols []string, values... any) error {
	if !isTableAllowed(tableName) {
		return fmt.Errorf("table %s not allowed", tableName)
	}
	if !areColumnsAllowed(tableName, cols) {
		return fmt.Errorf("columns %v not allowed for table %s", cols, tableName)
	}
	sqlQuery := fmt.Sprintf(`INSERT INTO %s (`, tableName)
	for i, v := range cols {
		if i != len(cols) - 1 { 
			sqlQuery = fmt.Sprint(sqlQuery, v, ", ") 
		} else {
			sqlQuery = fmt.Sprint(sqlQuery, v, ") VALUES (")
		}
	}
	
	for i := range values {
		if i != len(values) - 1 { 
			sqlQuery = fmt.Sprint(sqlQuery, "$", i + 1, ", ") 
		} else {
			sqlQuery = fmt.Sprint(sqlQuery, "$", i +1 , ");")
		}
	}
	_, err := d.pool.Exec(ctx, sqlQuery, values...)
	if err != nil {
		return err
	}	
	return nil 
}
