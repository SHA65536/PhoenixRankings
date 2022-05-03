CREATE TABLE IF NOT EXISTS "leaderboards" (
	"dbid"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"timestamp"	INTEGER NOT NULL,
	"id"	INTEGER NOT NULL,
	"name"	TEXT NOT NULL,
	"level"	INTEGER NOT NULL,
	"exp"	INTEGER NOT NULL,
	"fame"	INTEGER NOT NULL,
	"job"	INTEGER NOT NULL,
	"image"	TEXT,
	"restriction"	INTEGER NOT NULL
);
CREATE TABLE IF NOT EXISTS "players" (
	"id"	INTEGER NOT NULL PRIMARY KEY,
	"name"	TEXT NOT NULL
);