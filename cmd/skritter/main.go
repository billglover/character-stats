package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/billglover/character-stats/database"
	"github.com/billglover/character-stats/skritter"
)

func main() {

	log := log.New(os.Stdout, "api : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	err := run(log)
	if err != nil {
		log.Println("main: error:", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {

	var token string
	flag.StringVar(&token, "token", "", "Skritter API token")

	var dbFile string
	flag.StringVar(&dbFile, "db", "skritter.db", "path to the SQLite DB")

	flag.Parse()

	if token == "" {
		return errors.New("must provide -token")
	}

	if dbFile == "" {
		return errors.New("must specify the location of the db with -db")
	}

	//==========
	// Database
	log.Println("begin connecting to database")
	db, err := database.Open(database.Config{
		DSN: dbFile,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		log.Println("stopping database support")
		db.Close()
	}()

	err = database.StatusCheck(context.TODO(), db)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	log.Println("starting database migrations")
	err = database.Migrate(context.TODO(), db)
	if err != nil {
		return fmt.Errorf("migrating db: %w", err)
	}
	log.Println("completed database migrations")

	//=======
	// Skritter
	client := skritter.NewClient(token)
	vs, err := client.Vocab()
	if err != nil {
		return err
	}

	type Vocab struct {
		ID               string    `db:"id"`
		Lang             string    `db:"lang"`
		Priority         int       `db:"priority"`
		Style            string    `db:"style"`
		Audio            string    `db:"audio_url"`
		Toughness        int       `db:"toughness"`
		HeisigDefinition string    `db:"heisig_definition"`
		Ilk              string    `db:"ilk"`
		Writing          string    `db:"writing"`
		ToughnessString  string    `db:"toughness_string"`
		Definition       string    `db:"definition_en"`
		Starred          bool      `db:"starred"`
		Reading          string    `db:"reading"`
		Created          time.Time `db:"created_at"`
	}

	ts := time.Now()

	vocab := make([]Vocab, len(vs))
	for i, v := range vs {
		en := v.Definitions["en"]

		vocab[i] = Vocab{
			ID:               v.ID,
			Lang:             v.Lang,
			Priority:         v.Priority,
			Style:            v.Style,
			Audio:            v.Audio,
			Toughness:        v.Toughness,
			HeisigDefinition: v.HeisigDefinition,
			Ilk:              v.Ilk,
			Writing:          v.Writing,
			ToughnessString:  v.ToughnessString,
			Definition:       en,
			Starred:          false, // Ommitted as source type is inconsistent between bool and int
			Reading:          v.Reading,
			Created:          ts,
		}
	}

	for v := range vocab {
		fmt.Println(vocab[v])
	}

	res, err := db.NamedExec(`INSERT INTO vocab (id,lang,priority,style,audio_url,toughness,heisig_definition,ilk,writing,toughness_string,definition_en,starred,reading,created_at)
        VALUES (:id, :lang, :priority, :style, :audio_url, :toughness, :heisig_definition, :ilk, :writing, :toughness_string, :definition_en, :starred, :reading, :created_at)
		ON CONFLICT(id) DO UPDATE SET created_at=excluded.created_at;`, vocab)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Println("Inserted:", n)

	return nil
}
