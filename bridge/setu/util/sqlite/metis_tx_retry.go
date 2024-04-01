package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrMetisTxRetryNotFound = errors.New("metis tx retry not found")

type MetisTxRetry struct {
	ID     uint64 `json:"id" sql:"id"`
	TxData string `json:"tx_data" sql:"tx_data"`
}

type MetisTxRetryTable struct {
	db  *sql.DB
	rdb *sql.DB
}

func InitMetisTxRetryTable(db *sql.DB) error {
	var createTable = `CREATE TABLE IF NOT EXISTS metis_tx_retry_table(id INTEGER PRIMARY KEY AUTOINCREMENT, tx_data LONGTEXT);`
	_, err := db.Exec(createTable)
	return err
}

func NewMetisTxRetryTable(db, rdb *sql.DB) *MetisTxRetryTable {
	return &MetisTxRetryTable{
		db:  db,
		rdb: rdb,
	}
}

func (c *MetisTxRetryTable) Insert(txData string) (int, error) {
	res, err := c.db.Exec("INSERT INTO metis_tx_retry_table(tx_data) VALUES(?);", txData)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *MetisTxRetryTable) Delete(id uint64) error {
	_, err := c.db.Exec("DELETE FROM metis_tx_retry_table where id=?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *MetisTxRetryTable) GetAllWaitPushMetisTxRetrys(limit, offset int64) ([]*MetisTxRetry, error) {
	rows, err := c.rdb.Query("SELECT * FROM metis_tx_retry_table ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allMetisTxRetrys []*MetisTxRetry
	for rows.Next() {
		var event MetisTxRetry
		if err = rows.Scan(&event.ID, &event.TxData); err != nil {
			return nil, err
		}
		allMetisTxRetrys = append(allMetisTxRetrys, &event)
	}

	return allMetisTxRetrys, nil
}
