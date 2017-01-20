package helper

type Helper interface {
	GetTablesQuery() string
	TruncateTablesCommand(tableNames []string) string
}
