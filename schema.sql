CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    oauth BOOLEAN NOT NULL,
    hash TEXT NOT NULL,
    displayName TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS verify (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL,
    hash TEXT NOT NULL,
    displayName TEXT NOT NULL,
    code INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    postId TEXT NOT NULL,
    email TEXT NOT NULL,
    ip TEXT NOT NULL,
    displayName TEXT NOT NULL,
    content TEXT NOT NULL,
    imageURL TEXT NOT NULL,
    userWhere TEXT NOT NULL,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    address TEXT NOT NULL,
    constructionName TEXT NOT NULL,
    roadName TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    postId TEXT NOT NULL,
    email TEXT NOT NULL,
    ip TEXT NOT NULL,
    displayName TEXT NOT NULL,
    content TEXT NOT NULL,
    imageURL TEXT NOT NULL,
    userWhere TEXT NOT NULL,
    latitude REAL NOT NULL,
    longitude REAL NOT NULL,
    address TEXT NOT NULL,
    constructionName TEXT NOT NULL,
    roadName TEXT NOT NULL
);
