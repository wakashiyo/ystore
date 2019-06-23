package ystore

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	dbs struct {
		sync.RWMutex
		dbs map[string]*Db
	}
	ErrKeyNotFound = errors.New("Error: key not found")
	mutex          = &sync.RWMutex{}
)

//
// DB : key value store(DB)の情報
// 	params
// - name 				: データベース名？
// - fk 					: ？
// - fv						: ?
// - keys					: データに紐づくキー
// - vals					: 実際に扱うデータを入れる箱？
// - cancelSyncer : ?
// - storemode		: ?
type Db struct {
	sync.RWMutex
	name         string
	fk           *os.File
	fv           *os.File
	keys         [][]byte
	vals         map[string]*Cmd
	cancelSyncer context.CancelFunc
	storemode    int
}

// Cmd : コマンドの情報??
// 	params
// - Seek 				: ??
// - Size					: ??
// - KeySeek			: ??
// - Val					: ??
type Cmd struct {
	Seek    uint32
	Size    uint32
	KeySeek uint32
	Val     []byte
}

// Config : 設定の情報？
// 	params
// - FileMode 				: ??
// - DirMode 					: ??
// - SyncInterval 		: ??
// - StoreMode				: ??
type Config struct {
	FileMode     int
	DirMode      int
	SyncInterval int
	StoreMode    int
}

func init() {
	dbs.dbs = make(map[string]*Db)
}

func newDb(f string, cfg *Config) (*Db, error) {
	var err error

	db := new(Db)
	db.Lock()
	defer db.Unlock()
	db.name = f
	db.keys = make([][]byte, 0)
	db.vals = make(map[string]*Cmd)
	db.storemode = cfg.StoreMode

	if cfg.FileMode == 0 {
		cfg.FileMode = DefaultConfig.FileMode
	}
	if cfg.DirMode == 0 {
		cfg.DirMode = DefaultConfig.DirMode
	}
	if db.storemode == 2 && db.name == "" {
		return db, nil
	}

	_, err = os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			if filepath.Dir(f) != "." {
				err = os.MkdirAll(filepath.Dir(f), os.FileMode(cfg.DirMode))
				if err != nil {
					return nil, err
				}
			}
		} else {
			return nil, err
		}
	}
	db.fv, err = os.OpenFile(f, os.O_CREATE|os.O_RDWR, os.FileMode(cfg.FileMode))
	if err != nil {
		return nil, err
	}
	db.fk, err = os.OpenFile(f+".idx", os.O_CREATE|os.O_RDWR, os.FileMode(cfg.FileMode))
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	b, err := ioutil.ReadAll(db.fk)
	if err != nil {
		return nil, err
	}
	buf.Write(b)
	var readSeek uint32
	for buf.Len() > 0 {

	}
	if cfg.SyncInterval > 0 {
		db.backgroundManager(cfg.SyncInterval)
	}
	return db, err
}

func (db *Db) backgroundManager(interval int) {

}
