package util

import (
	"fmt"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

// GetBridgeDBInstance get sington object for bridge-db
func GetBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		fmt.Println("leveldb file path:", filePath)
		bridgeDB, _ = leveldb.OpenFile(filePath, nil)
	})

	return bridgeDB
}

// CloseBridgeDBInstance closes bridge-db instance
func CloseBridgeDBInstance() {
	bridgeDBCloseOnce.Do(func() {
		if bridgeDB != nil {
			bridgeDB.Close()
		}
	})
}
