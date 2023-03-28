package cache

import (
	_ "embed"
	"gorm.io/gorm"
	"os"
	"sync"
	"time"
)

//go:embed template/create_table.tpl
var createSql string

type dbCache struct {
	Options
	lock *sync.RWMutex
	db   *gorm.DB
}

func (d *dbCache) Set(key string, value []byte) (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	ldc := LeigDataCache{
		Key:        key,
		Value:      value,
		Timeout:    d.timeout,
		CreateTime: time.Now(),
	}
	tx := d.db.Create(&ldc)
	if tx.Error != nil {
		tx = d.db.Model(&ldc).Updates(LeigDataCache{
			Key:        ldc.Key,
			Value:      ldc.Value,
			Timeout:    ldc.Timeout,
			CreateTime: ldc.CreateTime,
		})
	}
	ok = tx.RowsAffected > 0
	return
}

func (d *dbCache) Get(key string) (value []byte) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	ldc := LeigDataCache{}
	tx := d.db.Where("`key` = ?", key).First(&ldc)
	if tx.RowsAffected > 0 {
		if time.Now().Sub(ldc.CreateTime) <= ldc.Timeout {
			value = ldc.Value
		}
		d.db.Delete(ldc, ldc.Id)
	}
	return
}

func NewDbCache(db *gorm.DB, opts ...Option) (Cache, error) {
	dos := defaultOptions()
	for _, apply := range opts {
		apply(&dos)
	}
	err := generate(db)
	if err != nil {
		return nil, err
	}
	return &dbCache{Options: dos, lock: &sync.RWMutex{}, db: db}, nil
}

type LeigDataCache struct {
	Id         int64         `gorm:"id"`
	Key        string        `gorm:"key"`
	Value      []byte        `gorm:"value"`
	Timeout    time.Duration `gorm:"timeout"`
	CreateTime time.Time     `gorm:"create_time"`
}

func generate(db *gorm.DB) error {
	replace := map[string]string{
		"TABLE_NAME": "leig_data_cache",
	}
	sql := os.Expand(createSql, func(s string) string {
		return replace[s]
	})
	tx := db.Exec(sql)
	return tx.Error
}
