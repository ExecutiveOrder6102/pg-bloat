package bloat

import "strings"

func ReindexIndexCommand(schema, index string, concurrent bool) string {
	if concurrent {
		return "REINDEX INDEX CONCURRENTLY " + quoteIdent(schema) + "." + quoteIdent(index)
	}
	return "REINDEX INDEX " + quoteIdent(schema) + "." + quoteIdent(index)
}

func ReindexTableCommand(schema, table string, concurrent bool) string {
	if concurrent {
		return "REINDEX TABLE CONCURRENTLY " + quoteIdent(schema) + "." + quoteIdent(table)
	}
	return "REINDEX TABLE " + quoteIdent(schema) + "." + quoteIdent(table)
}

func quoteIdent(value string) string {
	escaped := strings.ReplaceAll(value, `"`, `""`)
	return `"` + escaped + `"`
}
