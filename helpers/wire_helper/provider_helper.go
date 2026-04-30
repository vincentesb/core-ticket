package wire_helper

import (
	"core-ticket/constants"

	"github.com/jmoiron/sqlx"
)

/*
ProviderDBMain retrieves the main database connection from the provided map of database connections.

Parameters:
- db (map[string]*sqlx.DB): A map containing database connections where the key is a string and the value is a pointer to sqlx.DB.

Returns:
- (*sqlx.DB): The main database connection retrieved from the map based on the constant DBMain.

Example:

	db := map[string]*sqlx.DB{
		constants.DBMain: mainDBConnection,
	}
	mainDB := ProviderDBMain(db)
*/
func ProviderDBMain(db map[string]*sqlx.DB) *sqlx.DB {
	return db[constants.DBMain]
}
