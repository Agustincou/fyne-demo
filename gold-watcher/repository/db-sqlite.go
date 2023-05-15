package repository

import (
	"database/sql"
	"errors"
	"time"
)

type SQLiteRepository struct {
	Repository
	Conn *sql.DB
}

func NewSQLiteRespository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		Conn: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
		create table if not exists holdings(
			id integer primary key autoincrement,
			amount real not null,
			purchase_date integer not null,
			purchase_price integer not null);
	`

	_, err := r.Conn.Exec(query)

	return err
}

func (r *SQLiteRepository) InsertHolding(h Holding) (*Holding, error) {
	query := `insert into holdings (amount, purchase_date, purchase_price) values (?, ?, ?)`

	res, err := r.Conn.Exec(query, h.Amount, h.PurchaseDate.Unix(), h.PurchasePrice)
	if err != nil {
		return nil, err
	}

	if h.ID, err = res.LastInsertId(); err != nil {
		return nil, err
	}

	return &h, nil
}

func (r *SQLiteRepository) AllHoldings() ([]Holding, error) {
	query := `select id, amount, purchase_date, purchase_price from holdings order by purchase_date`

	rows, err := r.Conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []Holding
	for rows.Next() {
		var h Holding
		var unixTime int64
		err := rows.Scan(
			&h.ID,
			&h.Amount,
			&unixTime,
			&h.PurchasePrice,
		)
		if err != nil {
			return nil, err
		}
		h.PurchaseDate = time.Unix(unixTime, 0)
		all = append(all, h)
	}

	return all, nil
}

func (r *SQLiteRepository) GetHoldingByID(id int) (*Holding, error) {
	row := r.Conn.QueryRow(`select id, amount, purchase_date, purchase_price from holdings where id = ?`, id)

	var h Holding
	var unixTime int64
	if err := row.Scan(
		&h.ID,
		&h.Amount,
		&unixTime,
		&h.PurchasePrice,
	); err != nil {
		return nil, err
	}

	h.PurchaseDate = time.Unix(unixTime, 0)

	return &h, nil
}

func (r *SQLiteRepository) UpdateHolding(id int64, updated Holding) error {
	if id == 0 {
		return errors.New("invalid id to update")
	}

	stmt := "update holdings set amount = ?, purchase_date = ?, purchase_price = ? where id = ?"
	res, err := r.Conn.Exec(stmt, updated.Amount, updated.PurchaseDate.Unix(), updated.PurchasePrice, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errUpdateFailed
	}

	return nil
}

func (r *SQLiteRepository) DeleteHolding(id int64) error {
	res, err := r.Conn.Exec("delete from holdings where id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errDeleteFailed
	}

	return nil
}
