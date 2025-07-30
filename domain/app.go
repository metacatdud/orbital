package domain

import (
	"database/sql"
	"fmt"
	database "orbital/pkg/db"
)

type App struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Icon        string   `json:"icon"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Namespace   string   `json:"namespace"`
	OwnerKey    string   `json:"ownerKey"`
	OwnerURL    string   `json:"ownerUrl"`
	Labels      []string `json:"labels"`
	ParentID    string   `json:"parentId"`
	IsExternal  bool     `json:"isExternal"`
}

type appRow struct {
	ID          string
	Name        sql.NullString
	Icon        sql.NullString
	Version     sql.NullString
	Description sql.NullString
	Namespace   sql.NullString
	OwnerKey    sql.NullString
	OwnerURL    sql.NullString
	Labels      sql.NullString
	ParentID    sql.NullString
	IsExternal  sql.NullBool
}

type Apps []App

type AppRepository struct {
	db *database.DB
}

func NewAppRepository(db *database.DB) AppRepository {
	return AppRepository{db: db}
}

func (repo AppRepository) GetByID(id string) (*App, error) {
	query := `SELECT id, name, version, description, icon, namespace, owner_key, owner_url, labels, parent_id FROM applications WHERE id = ?`
	row := repo.db.Client().QueryRow(query, id)

	var appR appRow
	err := row.Scan(
		&appR.ID, &appR.Name, &appR.Version, &appR.Description,
		&appR.Icon, &appR.Namespace, &appR.OwnerKey, &appR.OwnerURL,
		&appR.Labels, &appR.ParentID, &appR.IsExternal,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	app := mapRowToApp(appR)

	return &app, nil
}

func (repo AppRepository) Find() (Apps, error) {
	query := `SELECT id, name, version, description, icon, namespace, owner_key, owner_url, labels, parent_id FROM applications`
	rows, err := repo.db.Client().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var apps Apps
	for rows.Next() {
		var appR appRow
		err = rows.Scan(
			&appR.ID, &appR.Name, &appR.Version, &appR.Description,
			&appR.Icon, &appR.Namespace, &appR.OwnerKey, &appR.OwnerURL,
			&appR.Labels, &appR.ParentID, &appR.IsExternal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		apps = append(apps, mapRowToApp(appR))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return apps, nil
}

func (repo AppRepository) FindOnlyStandalone() (Apps, error) {
	query := `SELECT id, name, version, description, icon, namespace, owner_key, owner_url, labels, parent_id, is_external FROM applications WHERE parent_id IS NULL`
	rows, err := repo.db.Client().Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var apps Apps
	for rows.Next() {
		var appR appRow
		err = rows.Scan(
			&appR.ID, &appR.Name, &appR.Version, &appR.Description,
			&appR.Icon, &appR.Namespace, &appR.OwnerKey, &appR.OwnerURL,
			&appR.Labels, &appR.ParentID, &appR.IsExternal,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		apps = append(apps, mapRowToApp(appR))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return apps, nil

}

func (repo AppRepository) FindByParentID(parentID string) (Apps, error) {
	query := `SELECT id, name, version, description, icon, namespace, owner_key, owner_url, labels, parent_id, is_external FROM applications WHERE parent_id = ?`
	rows, err := repo.db.Client().Query(query, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var apps Apps
	for rows.Next() {
		var appR appRow
		err = rows.Scan(
			&appR.ID, &appR.Name, &appR.Version, &appR.Description,
			&appR.Icon, &appR.Namespace, &appR.OwnerKey, &appR.OwnerURL,
			&appR.Labels, &appR.ParentID, &appR.IsExternal,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		apps = append(apps, mapRowToApp(appR))
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return apps, nil
}

func mapRowToApp(ar appRow) App {
	return App{
		ID:          ar.ID,
		Name:        nullToString(ar.Name),
		Icon:        nullToString(ar.Icon),
		Version:     nullToString(ar.Version),
		Description: nullToString(ar.Description),
		Namespace:   nullToString(ar.Namespace),
		OwnerKey:    nullToString(ar.OwnerKey),
		OwnerURL:    nullToString(ar.OwnerURL),
		Labels:      nullToStringSlice(ar.Labels),
		ParentID:    nullToString(ar.ParentID),
	}
}

func mapAppToRow(app App) appRow {
	return appRow{
		ID:          app.ID,
		Name:        stringToNull(app.Name),
		Icon:        stringToNull(app.Icon),
		Version:     stringToNull(app.Version),
		Description: stringToNull(app.Description),
		Namespace:   stringToNull(app.Namespace),
		OwnerKey:    stringToNull(app.OwnerKey),
		OwnerURL:    stringToNull(app.OwnerURL),
		Labels:      stringSliceToNull(app.Labels),
		ParentID:    stringToNull(app.ParentID),
	}
}
