package config

//InitSqls has the init sql for each table
var InitSqls = []string{
	`
    CREATE TABLE IF NOT EXISTS 'Article' (
	'id'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'name'	TEXT NOT NULL,
	'text'  TEXT NOT NULL,
	'chapter_id'	INTEGER,
	'law_id'	INTEGER,
	'reviewed' INTEGER
);`,
	`
    CREATE TABLE IF NOT EXISTS 'Chapter' (
	'id'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'name'	INTEGER NOT NULL,
	'title_id'	INTEGER,
	'law_id'	INTEGER,
	'reviewed' INTEGER
);`,
	`
    CREATE TABLE IF NOT EXISTS "Title" (
	'id'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'name'	TEXT,
	'law_id'	INTEGER,
	'reviewed' INTEGER
);`,
	`
    CREATE TABLE IF NOT EXISTS "Law" (
	'id'	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	'name'	TEXT NOT NULL,
	'approval_date'	TEXT NOT NULL,
	'publish_date'	TEXT NOT NULL,
	'journal'	TEXT,
	'intro'	TEXT,
	'reviewed' INTEGER,
	'revision' INTEGER
);`,
}
