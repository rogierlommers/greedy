package dao

import (
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	log "gopkg.in/inconshreveable/log15.v2"
)

const (
	sqlSelect        = "id, name, url, description, created"
	noDescription    = "greedy was unable to extract the meta description tag from your saved article."
	maxArticles      = 1000 // maximum of articles in database
	cleanupFrequency = 24   // in hours
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
		log.Crit("sqlite could not create database", "errormessage", err)
		return false
	}
	log.Info("sqlite succesfully initialized database file")
	return true
}

func check(e error) {
	if e != nil {
		log.Error("database interaction error", "message form db", e.Error())
	}
}

func Init(databasefile string) (db *sql.DB) {
	dbfileExists := false
	if _, err := os.Stat(databasefile); err == nil {
		dbfileExists = true
	}

	log.Debug("check if database file exists", "check result", dbfileExists)

	if true {
		var err error
		db, err = sql.Open("sqlite3", databasefile)
		if err != nil {
			log.Crit("error opening database file", "message", err)
		}
		log.Info("number of records in database", "amount", GetNumberOfRecords(db))
	}

	if !dbfileExists {
		initializeDatabase(db)
	}
	return
}

func Cleanup(db *sql.DB) {
	for {
		// query to run: DELETE FROM articles WHERE ROWID IN (SELECT ROWID FROM articles ORDER BY ROWID DESC LIMIT -1 OFFSET 10)
		log.Info("cleanup database file", "amount deleted", 1234)
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

// returns last x articles from db. If amount is 0, then all articles are returned
func GetLastArticles(db *sql.DB, amount int) (articleList []ArticleStruct) {
	var (
		stmt  *sql.Stmt
		rows  *sql.Rows
		limit string
	)

	if amount != 0 {
		limit = " LIMIT " + strconv.Itoa(amount)
	}

	tx, err := db.Begin()
	check(err)

	stmt, err = tx.Prepare("SELECT " + sqlSelect + " FROM articles ORDER BY id DESC" + limit)
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
		log.Error("error fetching article from database", "article id", id, "message", err)
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
	log.Info("start scraping article", "url", storedArticle.Url.String)

	if err != nil {
		log.Error("error while scraping article", "article id", id, "message", err)
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
	log.Debug("scrape information", "title", storedArticle.Name.String, "title length", len(storedArticle.Name.String))
	log.Debug("scrape information", "description", storedArticle.Description.String, "description length", len(storedArticle.Description.String))

	// after succesfull scraping, add page title (and more?) to article in db
	updateArticle(db, storedArticle)

	elapsed := time.Since(start)
	log.Info("scrape information", "stored article id", storedArticle.ID.Int64, "time elapsed", elapsed.String())
}
