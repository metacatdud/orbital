CREATE TABLE applications_old (
  id          TEXT PRIMARY KEY,
  name        TEXT NOT NULL,
  version     TEXT NOT NULL, -- revert to NOT NULL
  description TEXT,
  icon        TEXT NOT NULL DEFAULT 'fa-solid fa-satellite',
  namespace   TEXT NOT NULL, -- revert to NOT NULL
  is_external BOOLEAN NOT NULL DEFAULT 0,
  owner_key   TEXT,
  owner_url   TEXT,
  labels      TEXT,
  parent_id   TEXT DEFAULT NULL
);

INSERT INTO applications_old (
    id, name, version, description, icon, namespace,
    is_external, owner_key, owner_url, labels, parent_id
)
SELECT
    id, name, IFNULL(version, '-'), description, icon, IFNULL(namespace, ''),
    is_external, owner_key, owner_url, labels, parent_id
FROM applications;

DROP TABLE applications;
ALTER TABLE applications_old RENAME TO applications;
