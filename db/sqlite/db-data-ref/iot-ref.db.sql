BEGIN TRANSACTION;
DROP TABLE IF EXISTS "User";
CREATE TABLE IF NOT EXISTS "User" (
"id"	INTEGER PRIMARY KEY AUTOINCREMENT,
	"Username"	TEXT,
	"Password"	TEXT,
	"Salt"	INTEGER,
	"Active"	INTEGER,
	"Timestamp"	INTEGER
);
COMMIT;