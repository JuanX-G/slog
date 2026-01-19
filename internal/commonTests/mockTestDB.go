package commontests 

import (
	"fmt"
	"context"
	"slices"
)

type table struct {
	header []string
	rows [][]any
}

type MockDB struct {
    Tables map[string]table
}

func (m *MockDB) SelectAll(ctx context.Context, tableName string) ([][]any, error) {
	table, f := m.Tables[tableName]
	if !f {return nil, fmt.Errorf("no such table")}
	return table.rows, nil
}

func(m *MockDB) QueryForRow(ctx context.Context, tableName string, col string, value any) ([]any, error) {
	table, f := m.Tables[tableName]
	if !f {return nil, fmt.Errorf("no such table")}
	colIdx := -1
	for i, h := range table.header {
		if h == col {
			colIdx = i
			break
		}
	}
	if colIdx == -1 {return nil, fmt.Errorf("no such colums as %s", col)}

	for _, row := range table.rows {
		if len(row) < colIdx {return nil, fmt.Errorf("formatting error in mock db colum %s", col)}
		if row[colIdx] == value {return row, nil}
	}
	return nil, fmt.Errorf("no rows found where %s same as %s", col, value)
}

func (m *MockDB) CountWhere(ctx context.Context, tableName string, cols []string, values... any) (int32, error) {
	table, f := m.Tables[tableName]
	if !f {return -1, fmt.Errorf("no such table")}
	colIdxs := []int{}
	for i, h := range table.header {
		for _, col := range cols {
			if h == col {
				colIdxs = append(colIdxs, i)
			}
		}
	}
	if len(colIdxs) == 0 {return -1, fmt.Errorf("no such colums as %s", cols)}

	runningCount := 0
	for _, row := range table.rows {
		for _, idx := range colIdxs {
			if len(row) < idx {return -1, fmt.Errorf("formatting error in mock db table %s", tableName)}
			for _, val := range values {
				if row[idx] == val {runningCount++}
			}
		}
	}
	return int32(runningCount), nil
}

func (m *MockDB) QueryCountOffset(ctx context.Context, count int, offset int, table string, column string, value any) ([][]any, error) { 
	cnt := 0
	outRows := [][]any{}
	if offset > len(m.Tables[table].rows) {fmt.Errorf("to big of an offset")}
	for i := offset; i < len(m.Tables[table].rows); i++ {
		if cnt <= count {
			outRows = append(outRows, m.Tables[table].rows[i])
		} else {break}
	}
	return outRows, nil
}

func (m *MockDB) DeleteWhere(ctx context.Context, tableName string, cols []string, values... any) error {
	table, f := m.Tables[tableName]
	if !f {return fmt.Errorf("no such table")}

	colIdxs := []int{}
	for i, h := range table.header {
		for _, col := range cols {
			if h == col {
				colIdxs = append(colIdxs, i)
			}
		}
	}
	if len(colIdxs) == 0 {return fmt.Errorf("no such colums as %s", cols)}

	idxsForDeletion := []int{}
	for ridx, row := range table.rows {
		for _, idx := range colIdxs {
			if len(row) < idx {return fmt.Errorf("formatting error in mock db table %s", tableName)}
			for _, val := range values {
				if row[idx] == val {idxsForDeletion = append(idxsForDeletion, ridx)}
			}
		}
	}
	for _, idx := range idxsForDeletion {
		table.rows = slices.Delete(table.rows, idx, idx)
	}
	return nil 
}

// TODO: add emulating full insert... ehh...
func (m *MockDB) InsertInto(ctx context.Context, tableName string, cols []string, values... any) error {
	table, f := m.Tables[tableName]
	if !f {return fmt.Errorf("no such table")}
	table.rows = append(table.rows, values)
	return nil
}
