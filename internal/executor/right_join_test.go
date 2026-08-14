package executor_test

import (
	"testing"

	csvpkg "csvsql/internal/csv"
	"csvsql/internal/executor"
	"csvsql/internal/lexer"
	"csvsql/internal/parser"
	"csvsql/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func runSQLWithTables(t *testing.T, sql string, tables ...*csvpkg.Table) *executor.QueryResult {
	t.Helper()

	tokens, err := lexer.New(sql).AllTokens()
	require.NoError(t, err)
	stmt, err := parser.New(tokens).Parse()
	require.NoError(t, err)
	result, err := executor.New(tables).Execute(stmt)
	require.NoError(t, err)
	return result
}

func rightJoinTables() (*csvpkg.Table, *csvpkg.Table) {
	left := &csvpkg.Table{
		Name: "left_rows",
		Columns: []csvpkg.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "left_value", Type: types.TypeText},
		},
		Rows: [][]types.Value{
			{types.IntValue(1), types.TextValue("L1a")},
			{types.IntValue(1), types.TextValue("L1b")},
			{types.IntValue(2), types.TextValue("L2")},
			{types.IntValue(4), types.TextValue("L4")},
		},
	}
	right := &csvpkg.Table{
		Name: "right_rows",
		Columns: []csvpkg.Column{
			{Name: "id", Type: types.TypeInt},
			{Name: "right_value", Type: types.TypeText},
		},
		Rows: [][]types.Value{
			{types.IntValue(1), types.TextValue("R1a")},
			{types.IntValue(1), types.TextValue("R1b")},
			{types.IntValue(3), types.TextValue("R3")},
			{types.IntValue(5), types.TextValue("R5")},
		},
	}
	return left, right
}

func resultTexts(result *executor.QueryResult) [][]string {
	rows := make([][]string, len(result.Rows))
	for i, row := range result.Rows {
		rows[i] = make([]string, len(row))
		for j, value := range row {
			rows[i][j] = value.AsText()
		}
	}
	return rows
}

func TestRightJoinPreservesRightRowsAndNullExtendsLeft(t *testing.T) {
	left, right := rightJoinTables()
	result := runSQLWithTables(t, `
SELECT l.id AS left_id, l.left_value, r.id AS right_id, r.right_value
FROM left_rows l
RIGHT JOIN right_rows r ON l.id = r.id
ORDER BY r.id, l.left_value, r.right_value`, left, right)

	assert.Equal(t, [][]string{
		{"1", "L1a", "1", "R1a"},
		{"1", "L1a", "1", "R1b"},
		{"1", "L1b", "1", "R1a"},
		{"1", "L1b", "1", "R1b"},
		{"NULL", "NULL", "3", "R3"},
		{"NULL", "NULL", "5", "R5"},
	}, resultTexts(result))
	for _, row := range result.Rows[4:] {
		assert.True(t, row[0].IsNull(), "unmatched right rows must NULL-extend left columns")
		assert.True(t, row[1].IsNull(), "unmatched right rows must NULL-extend left columns")
	}
}
func TestRightJoinManyToManyMatches(t *testing.T) {
	left, right := rightJoinTables()
	result := runSQLWithTables(t, `
SELECT l.left_value, r.right_value
FROM left_rows l
RIGHT JOIN right_rows r ON l.id = r.id
WHERE r.id = 1
ORDER BY l.left_value, r.right_value`, left, right)

	assert.Equal(t, [][]string{
		{"L1a", "R1a"},
		{"L1a", "R1b"},
		{"L1b", "R1a"},
		{"L1b", "R1b"},
	}, resultTexts(result))
}

func TestRightJoinWhereSeesUnmatchedRightRows(t *testing.T) {
	left, right := rightJoinTables()
	result := runSQLWithTables(t, `
SELECT r.id, r.right_value
FROM left_rows l
RIGHT JOIN right_rows r ON l.id = r.id
WHERE l.id IS NULL
ORDER BY r.id`, left, right)

	assert.Equal(t, [][]string{{"3", "R3"}, {"5", "R5"}}, resultTexts(result))
}

func TestRightJoinAggregationIncludesUnmatchedRightRows(t *testing.T) {
	left, right := rightJoinTables()
	result := runSQLWithTables(t, `
SELECT r.id, COUNT(l.left_value) AS left_count
FROM left_rows l
RIGHT JOIN right_rows r ON l.id = r.id
GROUP BY r.id
ORDER BY r.id`, left, right)

	assert.Equal(t, [][]string{{"1", "4"}, {"3", "0"}, {"5", "0"}}, resultTexts(result))
}
