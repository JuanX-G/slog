package database

import (
	"context"
	"fmt"

)

func (d *DB) SelectAll(ctx context.Context, tableName string) ([][]any, error) {
	if !isTableAllowed(tableName) {
		return nil, fmt.Errorf("table %s not allowed", tableName)
	}

	sqlQuery := fmt.Sprintf(`SELECT * FROM %s`, tableName)
	rows, err := d.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rVals [][]any
	for (true) {
		if f := rows.Next(); !f {
			break;
		}
		aVals, err := rows.Values()
		if err != nil {
			return nil, err
		}
		rVals = append(rVals, aVals)
	}
	if len(rVals) == 0 {
		return nil, fmt.Errorf("No such rows")
	}
	return rVals, nil

}


func(d *DB) QueryForRow(ctx context.Context, tableName string, col string, value any) ([]any, error) {
	if !isTableAllowed(tableName) {
		return nil, fmt.Errorf("table %s not allowed", tableName)
	}
	var cols = []string{col}
	if !areColumnsAllowed(tableName, cols) {
		return nil, fmt.Errorf("columns %v not allowed for table %s", cols, tableName)
	}
	sqlQuery := fmt.Sprintf(`SELECT * FROM %s WHERE %s= $1;`, tableName, col)
	rows, err := d.pool.Query(ctx, sqlQuery, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if f := rows.Next(); !f {
		return nil, NotFoundError{queryParam: value}
	}
	vals, err := rows.Values()
	if err != nil {
		return nil, err
	}

	return vals, nil
}

