package store

import (
	"git.netflux.io/rob/solar-toolkit/inverter"
	"github.com/jmoiron/sqlx"
)

type PostgresStore struct {
	db *sqlx.DB
}

func NewSQL(db *sqlx.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) InsertETRuntimeData(runtimeData *inverter.ETRuntimeData) error {
	return nil
}
