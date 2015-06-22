package dao

import (
	"database/sql"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

const (
	sqlSelect        = "id, name, url, description, created"
	noDescription    = "go-read was unable to extract the meta description tag from your saved article."
	maxArticles      = 1000 // maximum of articles in database
	cleanupFrequency = 1    // in hours
)

// articleStruct holds the data fetched from a single row
type ArticleStruct struct {
	ID          sql.NullInt64
	Name        sql.NullString
	Url         sql.NullString
	Description sql.NullString
	Created     time.Time
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

func check(e error) {
	if e != nil {
		glog.Error(e)
	}
}

func Init(databasefile string) (db *sql.DB) {
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

func Cleanup(db *sql.DB) {
	for {

		// query to run: DELETE FROM articles WHERE ROWID IN (SELECT ROWID FROM articles ORDER BY ROWID DESC LIMIT -1 OFFSET 10)

		glog.Infof("cleanum removed #articles --> %d", 1234)

		time.Sleep(cleanupFrequency * time.Hour)

	}
}

func GetNumberOfRecords(db *sql.DB) (numberOfArticles int) {
	err := db.QueryRow("SELECT COUNT(1) FROM articles").Scan(&numberOfArticles)
	check(err)
	return numberOfArticles
}

func SaveArticle(db *sql.DB, url string) (lastInsertID int64) {
	stmt, err := db.Prepare("INSERT INTO articles (url, name) VALUES (?, ?)")
	check(err)

	result, err := stmt.Exec(url, url)
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

func getArticleById(db *sql.DB, id int64) (updatedArticle ArticleStruct) {
	var err error
	var stmt *sql.Stmt

	stmt, err = db.Prepare("SELECT " + sqlSelect + " FROM articles WHERE id = ?")
	check(err)

	err = stmt.QueryRow(id).Scan(&updatedArticle.ID, &updatedArticle.Name, &updatedArticle.Url, &updatedArticle.Description, &updatedArticle.Created)
	if err != nil {
		glog.Error("no record found with id --> ", id, err)
	}
	return
}

func updateArticle(db *sql.DB, updatedArticle ArticleStruct) (rowsAffected int64) {
	stmt, err := db.Prepare("UPDATE articles SET name = ?, description = ? WHERE id = ?")
	check(err)

	result, err := stmt.Exec(updatedArticle.Name, updatedArticle.Description, updatedArticle.ID)
	check(err)

	rowsAffected, err = result.RowsAffected()
	check(err)
	return
}

func ScrapeArticle(db *sql.DB, id int64) {
	// time function duration
	start := time.Now()

	// storedArticle contains information stored in db which need to be updated through scraping
	storedArticle := getArticleById(db, id)

	// init goquery
	doc, err := goquery.NewDocument(storedArticle.Url.String)
	glog.Info("scraping article with url --> ", storedArticle.Url.String)
	if err != nil {
		glog.Error("error while scraping article with id %d -- > ", storedArticle.ID, err)
		return
	}

	// start scraping page title
	doc.Find("head").Each(func(i int, s *goquery.Selection) {
		pageTitle := s.Find("title").Text()
		storedArticle.Name = sql.NullString{String: strings.TrimSpace(pageTitle), Valid: true}
	})

	// now get meta description field
	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if name, _ := s.Attr("name"); strings.EqualFold(name, "description") {
			description, _ := s.Attr("content")
			storedArticle.Description = sql.NullString{String: strings.TrimSpace(description), Valid: true}
		}
	})

	// if unable to scrape title, then use url
	if len(storedArticle.Name.String) == 0 {
		storedArticle.Name.String = storedArticle.Url.String
	}

	// if unable to scrape description, then use default text
	if len(storedArticle.Description.String) == 0 {
		storedArticle.Description.String = noDescription
	}

	// debugging info
	glog.Infof("scraped title --> %s (length: %d)", storedArticle.Name.String, len(storedArticle.Name.String))
	glog.Infof("scraped description --> %s (length: %d)", storedArticle.Description.String, len(storedArticle.Description.String))

	// after succesfull scraping, add page title (and more?) to article in db
	updateArticle(db, storedArticle)

	elapsed := time.Since(start)
	glog.Infof("scraping article %d completed in %s", storedArticle.ID.Int64, elapsed.String())
}
