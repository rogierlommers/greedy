package dao

import (
	"database/sql"

	"github.com/golang/glog"
)

func Init(db *sql.DB) bool {
	glog.Info("initializing databse")
	initializeDatabase(db)
	return true
}

func initializeDatabase(db *sql.DB) bool {
	sqlStmt := `
  DROP TABLE IF EXISTS attributes;
	CREATE TABLE attributes (
		id          INTEGER NOT NULL PRIMARY KEY,
		name        TEXT,
		alias       TEXT,
		parent_id   INTEGER,
		frequency   INTEGER DEFAULT 0,
		mark        INTEGER DEFAULT 0,
		-- pwd      TEXT,

		value_text  TEXT,
		value_blob  BLOB,
		value_int   INTEGER,
		value_real  REAL,
		value_time  DATETIME,

		accessed_at DATETIME,
		updated_at  DATETIME,
		deleted_at  DATETIME,
		created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

  CREATE UNIQUE INDEX IF NOT EXISTS index_on_alias        ON attributes (alias);
  CREATE        INDEX IF NOT EXISTS index_on_name         ON attributes (name);
  CREATE        INDEX IF NOT EXISTS index_on_value_text   ON attributes (value_text);
  CREATE        INDEX IF NOT EXISTS index_on_value_blob   ON attributes (value_blob);
  CREATE        INDEX IF NOT EXISTS index_on_value_int    ON attributes (value_int);
  CREATE        INDEX IF NOT EXISTS index_on_value_real   ON attributes (value_real);
  CREATE        INDEX IF NOT EXISTS index_on_accessed_at  ON attributes (accessed_at);
  CREATE        INDEX IF NOT EXISTS index_on_deleted_at   ON attributes (deleted_at);
  CREATE        INDEX IF NOT EXISTS index_on_frequency    ON attributes (frequency);
  CREATE        INDEX IF NOT EXISTS index_on_mark         ON attributes (mark);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		glog.Fatal(err)
		return false
	}
	glog.Info("repository initiated")
	return true
}
