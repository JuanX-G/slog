package database

import (
	"context"
	"strconv"
	"fmt"
)

func (d *DB) DeleteWhere(ctx context.Context, tableName string, cols []string, values... any) error {
	if !isTableAllowed(tableName) {
		return fmt.Errorf("table %s not allowed", tableName)
	}
	if !areColumnsAllowed(tableName, cols) {
		return fmt.Errorf("columns %v not allowed for table %s", cols, tableName)
	}
	sqlQuery := fmt.Sprintf(`DELETE FROM %s `, tableName)
	sqlQuery = fmt.Sprint(sqlQuery, "WHERE ")
	for i, v := range cols {
		if i != len(cols) - 1 { 
			sqlQuery = fmt.Sprint(sqlQuery, v, " = $", strconv.Itoa(i + 1), " AND ") 
		} else {
			sqlQuery = fmt.Sprint(sqlQuery, v, " = $", strconv.Itoa(i + 1)) 
		}
	}
	_, err := d.pool.Exec(ctx, sqlQuery, values...)
	if err != nil {
		return err
	}	
	return nil 
}
