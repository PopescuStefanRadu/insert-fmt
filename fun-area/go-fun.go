package main

import (
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
	for _, inserts := range insertsByTable {
		if len(inserts) > 1 {
			for _, sameTable := range inserts {
				var orderedColumns = make(sqlparser.Columns, len(sameTable.Columns))
				copy(sameTable.Columns, orderedColumns)
				sort.Slice(orderedColumns, func(i, j int) bool {
					return orderedColumns[i].CompliantName() < orderedColumns[j].CompliantName()
				})
				colNames := make([]string, len(sameTable.Columns))
				for _, column := range orderedColumns {
					colNames = append(colNames, column.Lowered())
				}
				columnsAsString := strings.Join(colNames, ",")
				insertsByColumns[columnsAsString] = append(insertsByColumns[columnsAsString], sameTable)
			}
		}
	}
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
			_, _ = os.Stderr.WriteString("Non-insert statement present")
		}
	}
	return statements
}

func getInsertsByTable(inserts []*sqlparser.Insert) map[string][]*sqlparser.Insert {
	var insertsByTable = make(map[string][]*sqlparser.Insert)
	for _, insert := range inserts {
		insertsByTable[insert.Table.Name.String()] = append(insertsByTable[insert.Table.Name.String()], insert)
	}
	return insertsByTable
}

//if insert, ok := statement.(*sqlparser.Insert); ok {
//	statements = append(statements, insert)
//} else {
//}
