package models

const (
	CreateTableQuery       = `CREATE TABLE %s (id SERIAL PRIMARY KEY, x  FLOAT NOT NULL, y  FLOAT NOT NULL);` //ID автоинкрементируемый
	CheckTableQuery        = `SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1);`
	CheckDataInTablesQuery = `SELECT EXISTS (SELECT 1 FROM measurements LIMIT 1);`
	DeleteQuery            = `DELETE FROM measurements;`
	GetListOfTablesQuery   = `SELECT measurement FROM information_schema.tables WHERE table_schema = 'public';`
	IDquery                = `SELECT get_next_measurement_number($1, $2);`
	InsertQuery            = `INSERT INTO measurements (x, y) VALUES ($1, $2);`
	SelectDataQuery        = `SELECT id, x, y FROM measurements;`
	// RenameTableQuery       = `ALTER TABLE measurements RENAME COLUMN x TO $1, RENAME COLUMN y TO $2;`
)
