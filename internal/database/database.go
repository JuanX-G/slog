package database

import (
	"context"
)

/* db-interface for the repo-pattern */
type Database interface {
    ConfigureDB() error
    InsertInto(ctx context.Context, table string, cols []string, values ...any) error
    QueryForRow(ctx context.Context, table, col string, value any) ([]any, error)
    QueryCountOffset(ctx context.Context, count int, offset int, table string, col string, value any) ([][]any, error)
    CountWhere(ctx context.Context, table string, cols []string, values ...any) (int32, error)
    DeleteWhere(ctx context.Context, table string, cols []string, values ...any) error
    SelectAll(ctx context.Context, table string) ([][]any, error)
}

