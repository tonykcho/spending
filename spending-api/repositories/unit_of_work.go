package repositories

import "database/sql"

type UnitOfWork interface {
	BeginTx() (*sql.Tx, error)
	CommitOrRollback(tx *sql.Tx, err error)
}

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *unitOfWork {
	return &unitOfWork{db: db}
}

func (uow *unitOfWork) BeginTx() (*sql.Tx, error) {
	return uow.db.Begin()
}

func (uow *unitOfWork) CommitOrRollback(tx *sql.Tx, err error) {
	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
}
