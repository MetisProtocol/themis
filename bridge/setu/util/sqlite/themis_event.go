package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var ErrThemisEventNotFound = errors.New("themis event not found")

type ThemisEvent struct {
	ID        uint64 `json:"id" sql:"id"`
	EventName string `json:"event_name" sql:"event_name"`
	EventLog  string `json:"event_log" sql:"event_log"`
}

type ThemisEventTable struct {
	db  *sql.DB
	rdb *sql.DB
}

func InitThemisEventTable(db *sql.DB) error {
	var createTable = `CREATE TABLE IF NOT EXISTS themis_event_table(id INTEGER PRIMARY KEY AUTOINCREMENT, event_name VARCHAR(128) NOT NULL, event_log LONGTEXT);`
	_, err := db.Exec(createTable)
	return err
}

func NewThemisEventTable(db, rdb *sql.DB) *ThemisEventTable {
	return &ThemisEventTable{
		db:  db,
		rdb: rdb,
	}
}

func (c *ThemisEventTable) Insert(eventName, eventLog string) (int, error) {
	res, err := c.db.Exec("INSERT INTO themis_event_table(event_name,event_log) VALUES(?,?);", eventName, eventLog)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (c *ThemisEventTable) Delete(id uint64) error {
	_, err := c.db.Exec("DELETE FROM themis_event_table where id=?;", id)
	if err != nil {
		return err
	}
	return nil
}

func (c *ThemisEventTable) GetAllWaitPushThemisEventsByType(eventName string, limit, offset int64) ([]*ThemisEvent, error) {
	rows, err := c.rdb.Query("SELECT * FROM themis_event_table WHERE event_name=? ORDER BY id ASC LIMIT ? OFFSET ?", eventName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allThemisEvents []*ThemisEvent
	for rows.Next() {
		var event ThemisEvent
		if err = rows.Scan(&event.ID, &event.EventName, &event.EventLog); err != nil {
			return nil, err
		}
		allThemisEvents = append(allThemisEvents, &event)
	}

	return allThemisEvents, nil
}
