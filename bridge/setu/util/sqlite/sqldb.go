package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

const sqliteOpenMode = "file:%v?_busy_timeout=5000&_journal=WAL&_sync=NORMAL&cache=shared&mode=rwc"

type SqliteDB struct {
	BridgeSqliteMetisTx       *MetisTxTable
	BridgeSqliteMetisTxRetry  *MetisTxRetryTable
	BridgeSqliteEthereumEvent *EthereumEventTable
	BridgeSqliteThemisEvent   *ThemisEventTable
}

var SqliteClient SqliteDB
var bridgeSqliteDBOnce sync.Once
var bridgeSqliteDBCloseOnce sync.Once

// GetBridgeDBInstance get sington object for bridge-sqlite-db
func GetBridgeSqlDBInstance(filePath string) *SqliteDB {
	bridgeSqliteDBOnce.Do(func() {
		fmt.Println("sqlite file path:", filePath)
		if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(filePath, 0777)
		}

		metisTxFile := filepath.Join(filePath, "metis_tx")
		fmt.Println("metis_tx_cache sqlite file path:", metisTxFile)
		if _, err := os.Stat(metisTxFile); errors.Is(err, os.ErrNotExist) {
			os.Create(metisTxFile)
		}
		metisTxDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, metisTxFile))
		metisTxDB.SetMaxOpenConns(1)
		metisTxReadDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, metisTxFile))
		metisTxReadDB.SetMaxOpenConns(1)
		InitMetisTxTable(metisTxDB)
		SqliteClient.BridgeSqliteMetisTx = NewMetisTxTable(metisTxDB, metisTxReadDB)

		metisTxRetryFile := filepath.Join(filePath, "metis_tx_retry")
		fmt.Println("metis_tx_retry sqlite file path:", metisTxRetryFile)
		if _, err := os.Stat(metisTxRetryFile); errors.Is(err, os.ErrNotExist) {
			os.Create(metisTxRetryFile)
		}
		metisTxRetryDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, metisTxRetryFile))
		metisTxRetryDB.SetMaxOpenConns(1)
		metisTxRetryReadDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, metisTxRetryFile))
		metisTxRetryReadDB.SetMaxOpenConns(1)
		InitMetisTxRetryTable(metisTxRetryDB)
		SqliteClient.BridgeSqliteMetisTxRetry = NewMetisTxRetryTable(metisTxRetryDB, metisTxRetryReadDB)

		ethereumEventFile := filepath.Join(filePath, "ethereum_event")
		fmt.Println("ethereum_event sqlite file path:", ethereumEventFile)
		if _, err := os.Stat(ethereumEventFile); errors.Is(err, os.ErrNotExist) {
			os.Create(ethereumEventFile)
		}
		ethereumEventDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, ethereumEventFile))
		ethereumEventDB.SetMaxOpenConns(1)
		ethereumEventReadDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, ethereumEventFile))
		ethereumEventReadDB.SetMaxOpenConns(1)
		InitEthereumEventTable(ethereumEventDB)
		SqliteClient.BridgeSqliteEthereumEvent = NewEthereumEventTable(ethereumEventDB, ethereumEventReadDB)

		themisEventFile := filepath.Join(filePath, "themis_event")
		fmt.Println("themis_event sqlite file path:", themisEventFile)
		if _, err := os.Stat(themisEventFile); errors.Is(err, os.ErrNotExist) {
			os.Create(themisEventFile)
		}
		themisEventDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, themisEventFile))
		themisEventDB.SetMaxOpenConns(1)
		themisEventReadDB, _ := sql.Open("sqlite3", fmt.Sprintf(sqliteOpenMode, themisEventFile))
		themisEventReadDB.SetMaxOpenConns(1)
		InitThemisEventTable(themisEventDB)
		SqliteClient.BridgeSqliteThemisEvent = NewThemisEventTable(themisEventDB, themisEventReadDB)
	})
	return &SqliteClient
}

// CloseBridgeSqlDBInstance closes bridge-sqlite-db instance
func CloseBridgeSqlDBInstance() {
	bridgeSqliteDBCloseOnce.Do(func() {
		if SqliteClient.BridgeSqliteMetisTx.db != nil {
			SqliteClient.BridgeSqliteMetisTx.db.Close()
		}
		if SqliteClient.BridgeSqliteMetisTx.rdb != nil {
			SqliteClient.BridgeSqliteMetisTx.rdb.Close()
		}
		if SqliteClient.BridgeSqliteMetisTxRetry.db != nil {
			SqliteClient.BridgeSqliteMetisTxRetry.db.Close()
		}
		if SqliteClient.BridgeSqliteMetisTxRetry.rdb != nil {
			SqliteClient.BridgeSqliteMetisTxRetry.rdb.Close()
		}
		if SqliteClient.BridgeSqliteEthereumEvent.db != nil {
			SqliteClient.BridgeSqliteEthereumEvent.db.Close()
		}
		if SqliteClient.BridgeSqliteEthereumEvent.rdb != nil {
			SqliteClient.BridgeSqliteEthereumEvent.rdb.Close()
		}
		if SqliteClient.BridgeSqliteThemisEvent.rdb != nil {
			SqliteClient.BridgeSqliteThemisEvent.db.Close()
		}
		if SqliteClient.BridgeSqliteThemisEvent.rdb != nil {
			SqliteClient.BridgeSqliteThemisEvent.db.Close()
		}
	})
}
