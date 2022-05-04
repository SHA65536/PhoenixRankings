CREATE TABLE IF NOT EXISTS "leaderboards" (
	"timestamp"	INTEGER NOT NULL,
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"rank"	INTEGER NOT NULL,
	"level"	INTEGER NOT NULL,
	"exp"	INTEGER NOT NULL,
	"fame"	INTEGER NOT NULL,
	"job"	INTEGER NOT NULL,
	"image"	TEXT,
	"restriction"	INTEGER NOT NULL
);