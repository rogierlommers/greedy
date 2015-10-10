package dao

import (
	"encoding/json"
	"fmt"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/boltdb/bolt"
	"github.com/rogierlommers/greedy/internal/common"
)

var (
	db   *bolt.DB
	open bool
)

type Article struct {
	ID    string
	Url   string
	Added time.Time
}

func Open() (err error) {
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(common.Databasefile, 0600, config)
	if err != nil {
		log.Crit("error creating bolt database", "message", err)
	}
	log.Info("database succesfully created", "filename", common.Databasefile)
	open = true
	return nil
}

func Close() {
	open = false
	db.Close()
}

func (a *Article) Save() error {
	if !open {
		return fmt.Errorf("db must be opened before saving")
	}
	err := db.Update(func(tx *bolt.Tx) error {
		articles, err := tx.CreateBucketIfNotExists([]byte("articles"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %s", err)
		}
		enc, err := a.encode()
		if err != nil {
			return fmt.Errorf("could not encode article %s:", err)
		}

		err = articles.Put([]byte(a.ID), enc)
		return err
	})
	return err
}

func (a *Article) encode() ([]byte, error) {
	enc, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func decode(data []byte) (*Article, error) {
	var a *Article
	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}
	return a, nil
}
