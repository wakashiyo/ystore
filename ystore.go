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

type Cmd struct {
	Seek    uint32
	Size    uint32
	KeySeek uint32
	Val     []byte
}

type Config struct {
	FileMode     int
	DirMode      int
	SyncInterval int
	Storemode    int
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
	db.storemode = cfg.Storemode

	if cfg.FileMode == 0 {
		cfg.FileMode = DefaultConfig.FileMode
	}
	if cfg.DirMode == 0 {
		cfg.DirMode = DefaultConfig.DirMode
	}
	if db.storemode == 2 && db.name == "" {
		return db, nil
	}

	//ファイル存在チェック
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
}
