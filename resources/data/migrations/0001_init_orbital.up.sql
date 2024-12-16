CREATE TABLE users
(
    id     TEXT PRIMARY KEY,
    name   TEXT NOT NULL,
    pubkey TEXT NOT NULL,
    access TEXT NOT NULL
);

CREATE TABLE containers
(
    id           TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    network_name TEXT NOT NULL,
    image        TEXT NOT NULL,
    ports        TEXT NOT NULL, -- JSON or string-encoded representation
    volumes      TEXT NOT NULL, -- JSON or string-encoded representation
    env_vars     TEXT,          -- JSON or string-encoded representation
    labels       TEXT           -- JSON or string-encoded representation
);

CREATE TABLE machines
(
    id   TEXT PRIMARY KEY,
    name TEXT NOT NULL
);
