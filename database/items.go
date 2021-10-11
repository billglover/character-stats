package database

import "time"

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
