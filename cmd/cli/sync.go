package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/billglover/character-stats/database"
	"github.com/billglover/character-stats/skritter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.Flags().String("token", "", "Skritter API token")
	syncCmd.MarkFlagRequired("token")

	syncCmd.Flags().String("db", "skritter.db", "Path to the local database")

	syncCmd.Flags().Bool("full", false, "perform a full database sync")

}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync local database with Skirtter",
	Long:  `Pull the latest data from Skritter and persist in a local database.`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) error {

	full, err := cmd.Flags().GetBool("full")
	if err != nil {
		return err
	}

	//==========
	// Config
	token, err := cmd.Flags().GetString("token")
	if err != nil {
		return err
	}

	dbFile, err := cmd.Flags().GetString("db")
	if err != nil {
		return err
	}

	//==========
	// Database
	db, err := database.Open(database.Config{
		DSN: dbFile,
	})
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}
	defer func() {
		db.Close()
	}()

	err = database.StatusCheck(context.TODO(), db)
	if err != nil {
		return fmt.Errorf("connecting to db: %w", err)
	}

	err = database.Migrate(context.TODO(), db)
	if err != nil {
		return fmt.Errorf("migrating db: %w", err)
	}

	//=======
	// Skritter
	client := skritter.NewClient(token)

	var since *time.Time

	if !full {
		result := struct {
			Changed time.Time `db:"changed"`
		}{}

		err := db.Get(&result, "SELECT changed FROM items ORDER BY changed DESC LIMIT 1;")
		if err != nil {
			return err
		}

		since = &result.Changed
	}

	vs, is, err := client.Vocab(since)
	if err != nil {
		return err
	}

	ts := time.Now()
	vocab := make([]database.Vocab, len(vs))
	items := make([]database.Item, len(is))

	for n, v := range vs {
		en := v.Definitions["en"]

		vocab[n] = database.Vocab{
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

		items[n] = database.Item{
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

	if len(vocab) != 0 {
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
	}

	if len(items) != 0 {
		var nt int64 = 0
		batch := 100

		for i := 0; i < len(items); i += batch {
			j := i + batch
			if j > len(items) {
				j = len(items)
			}
			q := `REPLACE INTO items (id,lang,style,changed,last,successes,time_studied,interval,next,reviews,previous_interval,part,vocab_id,previous_success,created_at)
	VALUES (:id, :lang, :style, :changed, :last, :successes, :time_studied, :interval, :next, :reviews, :previous_interval, :part, :vocab_id, :previous_success, :created_at);`

			res, err := db.NamedExec(q, items[i:j])
			if err != nil {
				return err
			}

			n, err := res.RowsAffected()
			if err != nil {
				return err
			}
			nt = nt + n
		}

		fmt.Println("Inserted Items:", nt)
	}

	return nil
}
