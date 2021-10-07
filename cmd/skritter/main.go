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
	vs, is, err := client.Vocab()
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

	type Item struct {
		ID               string    `db:"id"`
		Lang             string    `db:"lang"`
		Style            string    `db:"style"`
		Changed          time.Time `db:"changed"`
		Last             time.Time `db:"last"`
		Created          time.Time `db:"created_at"`
		Successes        int       `db:"successes"`
		TimeStudied      int       `db:"time_studied"`
		Interval         int       `db:"interval"`
		Next             time.Time `db:"next"`
		Reviews          int       `db:"reviews"`
		PreviousInterval int       `db:"previous_interval"`
		Part             string    `db:"part"`
		VocabId          string    `db:"vocab_id"`
		PreviousSuccess  bool      `db:"previous_success"`
	}

	ts := time.Now()
	vocab := make([]Vocab, len(vs))
	items := make([]Item, len(is))

	for n, v := range vs {
		en := v.Definitions["en"]

		vocab[n] = Vocab{
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

	for n, i := range is {
		changed := time.Unix(int64(i.Changed), 0)
		last := time.Unix(int64(i.Last), 0)
		created := time.Unix(int64(i.Created), 0)
		next := time.Unix(int64(i.Next), 0)

		items[n] = Item{
			ID:               i.ID,
			Lang:             i.Lang,
			Style:            i.Style,
			Changed:          changed,
			Last:             last,
			Created:          created,
			Successes:        i.Successes,
			TimeStudied:      i.TimeStudied,
			Interval:         i.Interval,
			Next:             next,
			Reviews:          i.Reviews,
			PreviousInterval: i.PreviousInterval,
			Part:             i.Part,
			VocabId:          i.VocabIds[0],
			PreviousSuccess:  i.PreviousSuccess,
		}
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

	fmt.Println("Inserted Vocab:", n)

	var nt int64 = 0
	batch := 100

	for i := 0; i < len(items); i += batch {
		j := i + batch
		if j > len(items) {
			j = len(items)
		}
		q := `REPLACE INTO items (id,lang,style,changed,last,successes,time_studied,interval,next,reviews,previous_interval,part,vocab_id,previous_success,created_at)
	VALUES (:id, :lang, :style, :changed, :last, :successes, :time_studied, :interval, :next, :reviews, :previous_interval, :part, :vocab_id, :previous_success, :created_at);`

		res, err = db.NamedExec(q, items[i:j])
		if err != nil {
			return err
		}

		n, err = res.RowsAffected()
		if err != nil {
			return err
		}
		nt = nt + n
	}

	fmt.Println("Inserted Items:", nt)

	return nil
}
