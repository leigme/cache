package cache

import (
	_ "embed"
	"gorm.io/gorm"
	"os"
	"strings"
	"sync"
	"time"
)

const DuplicateEntryError = "Duplicate entry"

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
	dct := DataCacheTable{
		Key:        key,
		Value:      value,
		Timeout:    d.timeout,
		CreateTime: time.Now(),
	}
	tx := d.db.Create(&dct)
	if tx.Error != nil && strings.Contains(tx.Error.Error(), DuplicateEntryError) {
		tx = d.db.Model(&dct).Where("`key` = ?", key).Updates(DataCacheTable{
			Value:      dct.Value,
			Timeout:    dct.Timeout,
			CreateTime: dct.CreateTime,
		})
	}
	ok = tx.RowsAffected > 0
	return
}

func (d *dbCache) Get(key string) (value []byte) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	dct := DataCacheTable{}
	tx := d.db.Where("`key` = ?", key).First(&dct)
	if tx.RowsAffected > 0 {
		if time.Now().Sub(dct.CreateTime) <= dct.Timeout {
			value = dct.Value
		} else {
			d.db.Delete(dct, dct.Id)
		}
	}
	return
}

func (d *dbCache) Remove(key string) (ok bool) {
	d.lock.Lock()
	defer d.lock.Unlock()
	tx := d.db.Where("`key` = ?", key).Delete(&DataCacheTable{})
	ok = tx.Error == nil
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

type DataCacheTable struct {
	Id         int64         `gorm:"id"`
	Key        string        `gorm:"key"`
	Value      []byte        `gorm:"value"`
	Timeout    time.Duration `gorm:"timeout"`
	CreateTime time.Time     `gorm:"create_time"`
}

func generate(db *gorm.DB) error {
	replace := map[string]string{
		"TABLE_NAME": "data_cache_table",
	}
	sql := os.Expand(createSql, func(s string) string {
		return replace[s]
	})
	tx := db.Exec(sql)
	return tx.Error
}
