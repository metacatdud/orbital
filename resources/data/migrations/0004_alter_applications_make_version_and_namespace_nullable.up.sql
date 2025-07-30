CREATE TABLE applications_new (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  version     TEXT,
  description TEXT,
  icon        TEXT NOT NULL DEFAULT 'fa-solid fa-satellite',
  namespace   TEXT,
  is_external BOOLEAN NOT NULL DEFAULT 0,
  owner_key   TEXT,
  owner_url   TEXT,
  labels      TEXT,
  parent_id   TEXT DEFAULT NULL
);

INSERT INTO applications_new (
    id, name, version, description, icon, namespace,
    is_external, owner_key, owner_url, labels, parent_id
)
SELECT
    id, name, version, description, icon, namespace,
    is_external, owner_key, owner_url, labels, parent_id
FROM applications;

DROP TABLE applications;
ALTER TABLE applications_new RENAME TO applications;
