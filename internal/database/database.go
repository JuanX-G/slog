package database

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

/* db-interface for the repo-pattern */
type Database interface {
    ConfigureDB() error
    InsertInto(ctx context.Context, table string, cols []string, values ...any) error
    QueryForRow(ctx context.Context, table, col string, value any) ([]any, error)
    QueryCountOffset(ctx context.Context, count int, offset int, table string, col string, value any) ([][]any, error)
    CountWhere(ctx context.Context, table string, cols []string, values ...any) (int32, error)
    DeleteWhere(ctx context.Context, table string, cols []string, values ...any) error
}

func (d *DB) CountWhere(ctx context.Context, tableName string, cols []string, values... any) (int32, error) {
	if !isTableAllowed(tableName) {
		return 0, fmt.Errorf("table %s not allowed", tableName)
	}
	if !areColumnsAllowed(tableName, cols) {
		return 0, fmt.Errorf("columns %v not allowed for table %s", cols, tableName)
	}
	sqlQuery := fmt.Sprintf(`SELECT COUNT(*) FROM %s `, tableName)
	sqlQuery = fmt.Sprint(sqlQuery, "WHERE ")
	for i, v := range cols {
		if i != len(cols) - 1 { 
			sqlQuery = fmt.Sprint(sqlQuery, v, " = $", strconv.Itoa(i + 1), " AND ") 
		} else {
			sqlQuery = fmt.Sprint(sqlQuery, v, " = $", strconv.Itoa(i + 1)) 
		}
	}
	var count int32
	err := d.pool.QueryRow(ctx, sqlQuery, values...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil 
}

