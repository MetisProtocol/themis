package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrEthereumEventNotFound = errors.New("Ethereum event not found")

type EthereumEvent struct {
	ID        uint64 `json:"id" sql:"id"`
	EventName string `json:"event_name" sql:"event_name"`
	EventLog  string `json:"event_log" sql:"event_log"`
}

type EthereumEventTable struct {
	db  *sql.DB
	rdb *sql.DB
}

func InitEthereumEventTable(db *sql.DB) error {
	var createTable = `CREATE TABLE IF NOT EXISTS ethereum_event_table(id INTEGER PRIMARY KEY AUTOINCREMENT, event_name VARCHAR(128) NOT NULL, event_log LONGTEXT);`
	_, err := db.Exec(createTable)
	return err
}

func NewEthereumEventTable(db, rdb *sql.DB) *EthereumEventTable {
	return &EthereumEventTable{
		db:  db,
		rdb: rdb,
	}
}

func (c *EthereumEventTable) Insert(eventName, eventLog string) (int, error) {
	res, err := c.db.Exec("INSERT INTO ethereum_event_table(event_name,event_log) VALUES(?,?);", eventName, eventLog)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *EthereumEventTable) Delete(id uint64) error {
	_, err := c.db.Exec("DELETE FROM ethereum_event_table where id=?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *EthereumEventTable) GetAllWaitPushEthereumEvents(limit, offset int64) ([]*EthereumEvent, error) {
	rows, err := c.rdb.Query("SELECT * FROM ethereum_event_table ORDER BY id ASC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allEthereumEvents []*EthereumEvent
	for rows.Next() {
		var event EthereumEvent
		if err = rows.Scan(&event.ID, &event.EventName, &event.EventLog); err != nil {
			return nil, err
		}
		allEthereumEvents = append(allEthereumEvents, &event)
	}

	return allEthereumEvents, nil
}
