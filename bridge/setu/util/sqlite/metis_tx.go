package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrMetisTxNotFound = errors.New("metis tx not found")

type MetisTx struct {
	ID     uint64 `json:"id" sql:"id"`
	TxHash string `json:"tx_hash" sql:"tx_hash"`
	TxData string `json:"tx_data" sql:"tx_data"`
	Pushed bool   `json:"pushed" sql:"pushed"` // pushed to l2
	Mined  bool   `json:"mined" sql:"mined"`   // l2 mined
}

type MetisTxTable struct {
	db  *sql.DB
	rdb *sql.DB
}

func InitMetisTxTable(db *sql.DB) error {
	var create = `CREATE TABLE IF NOT EXISTS metis_tx_table(id INTEGER PRIMARY KEY AUTOINCREMENT, tx_hash VARCHAR(128) NOT NULL, tx_data LONGTEXT, pushed BOOLEAN,mined BOOLEAN);
		CREATE UNIQUE INDEX IF NOT EXISTS uidx_tx_hash ON metis_tx_table(tx_hash);`
	_, err := db.Exec(create)
	return err
}

func NewMetisTxTable(db, rdb *sql.DB) *MetisTxTable {
	return &MetisTxTable{
		db:  db,
		rdb: rdb,
	}
}

func (c *MetisTxTable) Insert(tx *MetisTx) (int, error) {
	res, err := c.db.Exec("INSERT INTO metis_tx_table(tx_hash,tx_data,pushed,mined) VALUES(?,?,?,?)", tx.TxHash, tx.TxData, tx.Pushed, tx.Mined)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *MetisTxTable) DeleteExpiredDataByID(minID uint64) error {
	_, err := c.db.Exec("DELETE FROM metis_tx_table where id<=?", minID)
	if err != nil {
		return err
	}
	return nil
}

func (c *MetisTxTable) DeleteByID(id uint64) error {
	_, err := c.db.Exec("DELETE FROM metis_tx_table where id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *MetisTxTable) DeleteByTxHash(txHash string) error {
	_, err := c.db.Exec("DELETE FROM metis_tx_table where tx_hash=?", txHash)
	if err != nil {
		return err
	}
	return nil
}

func (c *MetisTxTable) GetMetisTxByID(id uint64) (*MetisTx, error) {
	// Query DB row based on ID
	row := c.rdb.QueryRow("SELECT id,tx_hash,tx_data,pushed,mined FROM metis_tx_table WHERE id=?", id)

	// Parse row into Activity struct
	var metisTx MetisTx
	var err error
	if err = row.Scan(&metisTx.ID, &metisTx.TxHash, &metisTx.TxData, &metisTx.Pushed, &metisTx.Mined); err == sql.ErrNoRows {
		return nil, ErrMetisTxNotFound
	}
	return &metisTx, nil
}

func (c *MetisTxTable) GetMetisTxByTxHash(txHash string) (*MetisTx, error) {
	// Query DB row based on ID
	row := c.rdb.QueryRow("SELECT * FROM metis_tx_table WHERE tx_hash=?", txHash)

	// Parse row into Activity struct
	var metisTx MetisTx
	if err := row.Scan(&metisTx.ID, &metisTx.TxHash, &metisTx.TxData, &metisTx.Pushed, &metisTx.Mined); err == sql.ErrNoRows {
		return nil, ErrMetisTxNotFound
	}
	return &metisTx, nil
}

func (c *MetisTxTable) GetLatestOne() (*MetisTx, error) {
	row := c.db.QueryRow("SELECT * FROM metis_tx_table ORDER BY id desc limit 1")

	// Parse row into Activity struct
	var metisTx MetisTx
	if err := row.Scan(&metisTx.ID, &metisTx.TxHash, &metisTx.TxData, &metisTx.Pushed, &metisTx.Mined); err == sql.ErrNoRows {
		return nil, ErrMetisTxNotFound
	}
	return &metisTx, nil
}

func (c *MetisTxTable) GetAllWaitPushMetisTxs(limit, offset, startID int64) ([]*MetisTx, error) {
	rows, err := c.rdb.Query("SELECT * FROM metis_tx_table WHERE id>? ORDER BY id ASC LIMIT ? OFFSET ?", startID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allMetisTxs []*MetisTx
	for rows.Next() {
		// Parse row into Activity struct
		var metisTx MetisTx
		if err = rows.Scan(&metisTx.ID, &metisTx.TxHash, &metisTx.TxData, &metisTx.Pushed, &metisTx.Mined); err != nil {
			return nil, err
		}
		allMetisTxs = append(allMetisTxs, &metisTx)
	}

	return allMetisTxs, nil
}

func (c *MetisTxTable) GetAllWaitPushMetisTxsByStartID(startID int64) ([]*MetisTx, error) {
	rows, err := c.rdb.Query("SELECT * FROM metis_tx_table WHERE id > ? ", startID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allMetisTxs []*MetisTx
	for rows.Next() {
		// Parse row into Activity struct
		var metisTx MetisTx
		if err = rows.Scan(&metisTx.ID, &metisTx.TxHash, &metisTx.TxData, &metisTx.Pushed, &metisTx.Mined); err != nil {
			return nil, err
		}
		allMetisTxs = append(allMetisTxs, &metisTx)
	}

	return allMetisTxs, nil
}
