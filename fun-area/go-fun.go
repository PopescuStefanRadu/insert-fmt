package main

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	reader, err := os.Open("/home/stefanradupopescu/work/personal/go/insert-fmt/test.sql")
	if err != nil {
		println(err)
		return
	}
	defer reader.Close()

	inserts := getInserts(reader)

	insertsByTable := getInsertsByTable(inserts)

	insertsByColumns := make(map[string][]*sqlparser.Insert)
	for _, inserts := range *insertsByTable {
		for _, sameTableInsert := range inserts {
			columns := sameTableInsert.Columns
			var orderedColumns = make(sqlparser.Columns, len(columns))
			copy(orderedColumns, columns)
			sort.Slice(orderedColumns, func(i, j int) bool {
				return orderedColumns[i].CompliantName() < orderedColumns[j].CompliantName()
			})
			colNames := make([]string, len(columns))
			for i, column := range orderedColumns {
				colNames[i] = column.Lowered()
			}
			columnsAsString := strings.Join(colNames, ",")
			insertsByColumns[columnsAsString] = append(insertsByColumns[columnsAsString], sameTableInsert)
		}
	}
	fmt.Println(insertsByColumns)
}

func getInserts(reader io.Reader) []*sqlparser.Insert {
	tokens := sqlparser.NewTokenizer(reader)

	var statements []*sqlparser.Insert
	for {
		statement, err := sqlparser.ParseNext(tokens)
		if err == io.EOF {
			break
		}
		if insert, ok := statement.(*sqlparser.Insert); ok {
			statements = append(statements, insert)
		} else {
			_, _ = os.Stderr.WriteString("[Warning] - Non-insert statement present\r\n")
		}
	}
	return statements
}

func getInsertsByTable(inserts []*sqlparser.Insert) *map[string][]*sqlparser.Insert {
	var insertsByTable = make(map[string][]*sqlparser.Insert)
	for _, insert := range inserts {
		insertsByTable[insert.Table.Name.String()] = append(insertsByTable[insert.Table.Name.String()], insert)
	}
	return &insertsByTable
}
