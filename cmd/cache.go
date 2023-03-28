package main

import (
	"fmt"
	"github.com/leigme/cache"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/local?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	c, err := cache.NewDbCache(db)
	fmt.Println(c.Set("abc", []byte("abcdefg")))
	fmt.Println(string(c.Get("abc")))
}
