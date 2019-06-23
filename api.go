package ystore

var DefaultConfig = &Config{
	FileMode:     0666,
	DirMode:      0777,
	SyncInterval: 0,
	StoreMode:    0,
}

// Open : データベースを開いて、データベースのオブジェクトを返す
//       データベースが存在しない場合は、新しく作成して返す
//			 データベースが存在する場合は、そのオブジェクトを読んで返す
//       Config(cfg)がnilの場合は、DefaultConfigを使用する
func Open(f string, cfg *Config) (*Db, error) {
	if cfg == nil {
		cfg = DefaultConfig
	}
	dbs.RLock()
	//存在する場合は、DBオブジェクトを返す
	db, ok := dbs.dbs[f]
	if ok {
		dbs.RUnlock()
		return db, nil
	}
	dbs.RUnlock()

	//RESEARCH : sync.RWMutexのLockとRLockの違い

	dbs.Lock()
	//存在しない場合は、新たにDBオブジェクトを作成する
	db, err := newDb(f, cfg)
	if err == nil {
		dbs.dbs[f] = db
	}
	dbs.Unlock()
	return db, err
}
