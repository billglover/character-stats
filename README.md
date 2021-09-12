# character-stats

Identify the Chinese characters I have most difficulty remembering

## Proposal

I would like to know which Chinese characters are the ones I find most challenging to learn.

**Solution:** Provide a web-page that shows me my top 10 most difficult Chinese characters. Difficulty should be determined by the number of incorrect attempts to recall a character. Difficult should be weighted to favour recent incorrect recall attempts.

**Implementation:**
* Backend queries study data from Skritter daily
* Study data is cached in a database for statistical analysis
* Backend exposes an API to allow data to be displayed
* Front end displays a list of challenging characters

**Considerations:**
* Backend needs to authenticate to query Skritter
* Front end doesn't need to authenticate for read-only API
