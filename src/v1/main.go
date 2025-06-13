package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"

	"zero-downtime-training/src"
)

type AssetManagerV1 struct {
	db *sql.DB
}

func (v1 *AssetManagerV1) CreateAsset(asset *src.Asset) error {
	query := "INSERT INTO assets (id, name, source) VALUES (?, ?, ?)"
	_, err := v1.db.Exec(query, asset.ID, asset.Name, asset.Source)
	if err != nil {
		return err
	}
	return nil
}

func (v1 *AssetManagerV1) DeleteAsset(id string) error {
	query := "DELETE FROM assets WHERE id = ?"
	_, err := v1.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

func (v1 *AssetManagerV1) ReadAsset(id string) (*src.Asset, error) {
	query := "SELECT id, name, source FROM assets WHERE id = ?"
	row := v1.db.QueryRow(query, id)

	asset := &src.Asset{}
	if err := row.Scan(&asset.ID, &asset.Name, &asset.Source); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err // Other error
	}
	return asset, nil
}

func (v1 *AssetManagerV1) UpdateAssetSourceByID(id string, name string) error {
	query := "UPDATE assets SET source = ? WHERE id = ?"
	_, err := v1.db.Exec(query, name, id)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	dbHost := src.GetEnvOrPanic(src.EnvDBHost)
	dbPort := src.GetEnvOrPanic(src.EnvDBPort)
	username := src.GetEnvOrPanic(src.EnvDBUser)
	password := src.GetEnvOrPanic(src.EnvDBPassword)
	schema := src.GetEnvOrPanic(src.EnvDBSchema)
	port := src.GetEnvOrPanic(src.EnvServerPort)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", username, password, dbHost, dbPort, schema)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mngr := &AssetManagerV1{db: db}
	srv := src.NewServer(port, mngr)

	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
