package helper

// Helper provides a way for the cleaner to find out list of tables to work on,
// also how to perform truncate on them
type Helper interface {
	GetTablesQuery() string
	TruncateTablesCommand(tableNames []string) string
}
