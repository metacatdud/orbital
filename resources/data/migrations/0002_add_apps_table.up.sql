CREATE TABLE applications
(
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    version     TEXT NOT NULL,
    description TEXT,
    icon        TEXT NOT NULL DEFAULT 'fa-solid fa-satellite',
    namespace   TEXT NOT NULL,
    owner_key   TEXT, -- Owner public key
    owner_url   TEXT,
    labels      TEXT, -- JSON or string-encoded representation
    parent_id   TEXT DEFAULT NULL
);
