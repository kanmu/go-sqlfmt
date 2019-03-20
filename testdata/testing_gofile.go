package sqlfmt

import (
	"database/sql"
)

func sendSQL() int {
	var id int
	var db *sql.DB
	db.QueryRow(`
select any ( select xxx from xxx ) from xxx where xxx limit xxx `).Scan(&id)
	return id
}
