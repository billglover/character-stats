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

-- Version: 0.2
-- Description: Create table items
CREATE TABLE IF NOT EXISTS "items" (
	"id" text,
	"lang" text,
	"style" text,
	"changed" datetime,
	"last" datetime,
	"successes" int,
	"time_studied" int,
	"interval" int,
	"next" datetime,
	"reviews" int,
	"previous_interval" int,
	"part" string,
	"vocab_id" string,
	"previous_success" boolean,
    "created_at" datetime,
	PRIMARY KEY (id)
);

CREATE INDEX "vocab_id" ON "items" ("vocab_id");
