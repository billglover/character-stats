package skritter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Client) Vocab() ([]Vocab, []Item, error) {

	vocab := map[string]Vocab{}
	items := map[string]Item{}

	path := c.BaseURL.String() + "/items"
	cursor := ""

	for {

		req, err := http.NewRequest(http.MethodPost, path, nil)
		if err != nil {
			return nil, nil, err
		}

		q := req.URL.Query()
		q.Add("include_vocabs", "true")
		q.Add("include_heisigs", "true")
		if cursor != "" {
			q.Add("cursor", cursor)
		}
		req.URL.RawQuery = q.Encode()

		req.Header.Add("Authorization", "Bearer "+c.token)
		req.Method = http.MethodGet

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, nil, err
		}
		defer resp.Body.Close()

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, nil, err
		}
		bodyString := string(bodyBytes)

		if resp.StatusCode != http.StatusOK {
			return nil, nil, fmt.Errorf("%s: %s", resp.Status, bodyString)
		}

		type Response struct {
			Cursor     string  `json:"cursor"`
			Items      []Item  `json:"Items"`
			Vocabs     []Vocab `json:"Vocabs"`
			StatusCode int     `json:"statusCode"`
		}

		var r Response
		err = json.Unmarshal(bodyBytes, &r)
		if err != nil {
			return nil, nil, err
		}

		cursor = r.Cursor

		for _, i := range r.Items {
			items[i.ID] = i
		}

		for _, v := range r.Vocabs {
			vocab[v.ID] = v
		}

		if cursor == "" {
			break
		}

		//time.Sleep(5 * time.Millisecond)
	}

	fmt.Println("Items retrieved:", len(items))
	is := make([]Item, len(items))
	n := 0
	for _, i := range items {
		is[n] = i
		n++
	}

	fmt.Println("Vocab retrieved:", len(vocab))
	vs := make([]Vocab, len(vocab))
	n = 0
	for _, v := range vocab {
		vs[n] = v
		n++
	}

	return vs, is, nil
}

// Item is the atomic unit of learning in Skritter. It is used to track
// reviews and learning. An item could be the writing for a character, or
// the tone of a word.
// Source Documentation: https://www.skritter.com/api/v0/docs/entities/items
type Item struct {
	Lang             string   `json:"lang"`
	Style            string   `json:"style"`
	Interval         int      `json:"interval"`
	Last             int      `json:"last"`
	Created          int      `json:"created"`
	Successes        int      `json:"successes"`
	TimeStudied      int      `json:"timeStudied"`
	Changed          int      `json:"changed"`
	Next             int      `json:"next"`
	Reviews          int      `json:"reviews"`
	PreviousInterval int      `json:"previousInterval,omitempty"`
	SectionIds       []string `json:"sectionIds"`
	VocabListIds     []string `json:"vocabListIds"`
	VocabIds         []string `json:"vocabIds"`
	Part             string   `json:"part"`
	PreviousSuccess  bool     `json:"previousSuccess"`
	ID               string   `json:"id"`
}

// Vocab is complementary to the Item entity. It and provides all the user-specific
// settings for a word, as well as all the information about the word.
// Source Documentation: https://www.skritter.com/api/v0/docs/entities/vocabs
type Vocab struct {
	Lang             string            `json:"lang"`
	Priority         int               `json:"priority"`
	Style            string            `json:"style"`
	Audio            string            `json:"audio"`
	Toughness        int               `json:"toughness"`
	DictionaryLinks  map[string]string `json:"dictionaryLinks"`
	HeisigDefinition string            `json:"heisigDefinition"`
	Ilk              string            `json:"ilk"`
	Writing          string            `json:"writing"`
	Audios           []Audio           `json:"audios"`
	AudioURL         string            `json:"audioURL"`
	ToughnessString  string            `json:"toughnessString"`
	Definitions      map[string]string `json:"definitions"`
	//Starred          bool              `json:"starred"` // Omitted as type varies between bool and int
	Reading string `json:"reading"`
	ID      string `json:"id"`
}

// Audio contains links to the pronounciation of a vocab.
type Audio struct {
	Source  string `json:"source"`
	Reading string `json:"reading"`
	ID      string `json:"id"`
	Writing string `json:"writing"`
	Mp3     string `json:"mp3"`
}
