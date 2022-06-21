package utils

import (
	"fmt"
	"github.com/advancemg/vimb-loader/pkg/models"
	"testing"
)

func TestSplitDbAndTable(t *testing.T) {
	db, table := SplitDbAndTable(models.DbBudgets)
	fmt.Println(db)
	fmt.Println(table)
}
