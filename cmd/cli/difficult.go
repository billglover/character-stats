package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/billglover/character-stats/database"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(difficultCmd)
	difficultCmd.Flags().IntP("number", "n", 3, "Number of characters to return")
	difficultCmd.Flags().String("db", "skritter.db", "Path to the local database")
	difficultCmd.Flags().Bool("tone", false, "include tone recall when determining difficulty")
	difficultCmd.Flags().Bool("reading", false, "include reading recall when determining difficulty")
	difficultCmd.Flags().Bool("writing", true, "include writing recall when determining difficulty")
	difficultCmd.Flags().Bool("definition", true, "include definition recall when determining difficulty")
}

var difficultCmd = &cobra.Command{
	Use:   "difficult",
	Short: "Identifies the most difficult caharacters",
	Long:  `Queries the local database to identify the most difficult characters.`,
	RunE:  difficult,
}

func difficult(cmd *cobra.Command, args []string) error {

	n, err := cmd.Flags().GetInt("number")
	if err != nil {
		return err
	}

	dbFile, err := cmd.Flags().GetString("db")
	if err != nil {
		return err
	}

	//==============================
	// Construct the Query Condition
	inclTone, _ := cmd.Flags().GetBool("tone")
	inclReading, _ := cmd.Flags().GetBool("reading")
	inclWriting, _ := cmd.Flags().GetBool("writing")
	inclDefinition, _ := cmd.Flags().GetBool("definition")

	var whereParts []string
	if !inclTone {
		whereParts = append(whereParts, `items.part != "tone"`)
	}
	if !inclReading {
		whereParts = append(whereParts, `items.part != "rdng"`)
	}
	if !inclWriting {
		whereParts = append(whereParts, `items.part != "rune"`)
	}
	if !inclDefinition {
		whereParts = append(whereParts, `items.part != "defn"`)
	}

	where := strings.Join(whereParts, " AND ")

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

	q := `
	SELECT
		vocab_id,
		part,
		reviews,
		successes,
		(successes*100/reviews) 'percent',
		writing,
		heisig_definition,
		definition_en,
		reading
	FROM
		items
	INNER JOIN vocab ON vocab_id = vocab.id
	WHERE
		` + where + `
	ORDER BY
		percent
	LIMIT :n;`

	params := struct {
		N int `db:"n"`
	}{
		N: n,
	}

	vocab := []struct {
		VocabID    string `db:"vocab_id"`
		Part       string `db:"part"`
		Writing    string `db:"writing"`
		Reviews    int    `db:"reviews"`
		Successes  int    `db:"successes"`
		Percent    int    `db:"percent"`
		HeisigDefn string `db:"heisig_definition"`
		Definition string `db:"definition_en"`
		Reading    string `db:"reading"`
	}{}

	err = database.NamedQuerySlice(context.TODO(), db, q, params, &vocab)
	if err != nil {
		return err
	}

	for _, i := range vocab {
		i.Definition = strings.ReplaceAll(i.Definition, "\n", " ")
		fmt.Printf("%02d%% %-8s %-2s %-15s %s %s\n", i.Percent, characterPart(i.Part), i.Writing, "("+i.Reading+")", "["+i.HeisigDefn+"]", i.Definition)
	}

	return nil
}

func characterPart(p string) string {
	switch p {
	case "defn":
		return "meaning"
	case "rune":
		return "writing"
	case "tone":
		return "tone"
	case "rdng":
		return "reading"
	default:
		return p
	}
}
