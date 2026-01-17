package database

import (
	"context"
	"strconv"
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

func (db *DB) QueryCountOffset(ctx context.Context, count int, offset int, table string, column string, value any) ([][]any, error) { 
	if !isTableAllowed(table) { 
		return nil, fmt.Errorf("table %s not allowed", table) 
	} 
	var cols = []string{column} 
	if !areColumnsAllowed(table, cols) { 
		return nil, fmt.Errorf("columns %v not allowed for table %s", cols, table) 
	} 
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1 ORDER BY date_created DESC LIMIT $2 OFFSET $3", table, column) 
	rows, err := db.pool.Query(ctx, query, value, count, offset)
	if err != nil { 
		return nil, err 
	} 
	defer rows.Close() 
	var result [][]any 
	for rows.Next() { 
		vals, err := rows.Values() 
		if err != nil { 
			return nil, err 
		} 
		result = append(result, vals) 
	} 
	if len(result) == 0 { 
		return nil, fmt.Errorf("no results found") 
	} 
	if count == -1 { 
		return result, nil 
	} 
	return result[:count], nil 
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

func(d *DB) QueryCount(ctx context.Context, count int, tableName, col string, value any) ([][]any, error){
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
	if len(rVals) < count {
		return nil, fmt.Errorf("Count too big")
	}
	if count == -1 {
		return rVals, nil
	}
	return rVals[:count + 1], nil
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
