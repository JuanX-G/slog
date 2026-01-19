package database_test

import (
	"fmt"
	common "slog-simple-blog/internal/commonTests"
	"testing"
)

func TestDatabase(t *testing.T) {
	mcdb := common.MockDB{}
	fmt.Print(mcdb)
}

