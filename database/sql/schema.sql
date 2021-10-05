-- Version: 0.1
-- Description: Create table vocab
CREATE TABLE IF NOT EXISTS "vocab" (
	"id" text,
	"lang" text,
	"priority" int,
	"style" text,
	"audio_url" text,
	"toughness" int,
	"heisig_definition" text,
	"ilk" text,
	"writing" text,
	"toughness_string" text,
	"definition_en" text,
	"starred" num,
	"reading" text,
    "created_at" datetime,
	PRIMARY KEY (id)
);