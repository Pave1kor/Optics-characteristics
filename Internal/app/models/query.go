package models

const (
	CreateTableQuery = `
	CREATE TABLE IF NOT EXISTS measurements (
    id TEXT NOT NULL,
    type TEXT NOT NULL,
    Date DATE NOT NULL,
    Number INTEGER NOT NULL,
    column_name TEXT NOT NULL,
    x DOUBLE PRECISION NOT NULL,
    y DOUBLE PRECISION NOT NULL,
    PRIMARY KEY (id, column_name)
);`

	DeleteQuery = `DELETE FROM measurements`

	GetListOfTablesQuery = `SELECT measurement FROM information_schema.tables WHERE table_schema = 'public';`

	IDquery = `SELECT get_next_measurement_number($1, $2)`

	InsertQuery = `INSERT INTO measurements (id, type, Date, Number, column_name, x, y) 
					VALUES ($1, $2, $3, $4, $5, $6, $7)`
)
