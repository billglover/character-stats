package database

import "time"

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
