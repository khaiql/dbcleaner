package helper

type Helper interface {
	GetTablesQuery() string
	TruncateTableCommand(tableName string) string
}
