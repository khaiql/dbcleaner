package instructor

type Instructor interface {
	GetTablesQuery() string
	TruncateTableCommand(tableName string) string
}
