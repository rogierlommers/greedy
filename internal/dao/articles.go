package dao

import (
	"database/sql"
	"os"
	"time"

	"github.com/golang/glog"
)

const sqlSelect = "id, name, url, description, created"

// articleStruct holds the data fetched from a row
type ArticleStruct struct {
	ID          sql.NullInt64
	Name        sql.NullString
	Url         sql.NullString
	Description sql.NullString
	Created     time.Time
}

func check(e error) {
	if e != nil {
		glog.Error(e)
	}
}

func Init(databasefile string) (db *sql.DB) {
	//	var db *sql.DB

	dbfileExists := false
	if _, err := os.Stat(databasefile); err == nil {
		dbfileExists = true
	}

	glog.Info("does dbfile exist? ", dbfileExists)

	if true {
		var err error
		db, err = sql.Open("sqlite3", databasefile)
		if err != nil {
			glog.Fatal(err)
		}

	}

	if !dbfileExists {
		initializeDatabase(db)
	}
	return
}

func initializeDatabase(db *sql.DB) bool {
	sqlStmt := `DROP TABLE IF EXISTS articles;
				CREATE TABLE articles (
					id          INTEGER NOT NULL PRIMARY KEY,
					name        TEXT,
					url         TEXT,
					description TEXT,
					created     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
				);`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		glog.Fatal(err)
		return false
	}
	glog.Info("repository initiated")
	return true
}

func SaveArticle(db *sql.DB, url string) (lastInsertID int64) {
	stmt, err := db.Prepare("INSERT INTO articles (url) VALUES (?)")
	check(err)

	result, err := stmt.Exec(url)
	check(err)

	lastInsertID, err = result.LastInsertId()

	check(err)
	return
}

func GetLastArticles(db *sql.DB) (articleList []ArticleStruct) {
	var stmt *sql.Stmt
	var rows *sql.Rows

	tx, err := db.Begin()
	check(err)

	stmt, err = tx.Prepare("SELECT " + sqlSelect + " FROM articles ORDER BY id DESC")
	check(err)
	defer stmt.Close()

	rows, err = stmt.Query()
	check(err)
	defer rows.Close()

	articleList = make([]ArticleStruct, 0, 0)

	for rows.Next() {
		article := ArticleStruct{}
		err = rows.Scan(&article.ID, &article.Name, &article.Url, &article.Description, &article.Created)
		check(err)
		articleList = append(articleList, article)
	}

	tx.Commit()
	return
}

func ScrapeArticle(id int64) {
	glog.Info("Scraping article with id --> ", id)
}
