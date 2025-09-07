package repositories

import "database/sql"

type UnitOfWork interface {
	WithTransaction(fn func(tx *sql.Tx) error) error
}

type unitOfWork struct {
	db *sql.DB
}

func NewUnitOfWork(db *sql.DB) *unitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) WithTransaction(fn func(tx *sql.Tx) error) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	return fn(tx)
}
