package db

import (
    "github.com/jinzhu/gorm"
    "log"
)

var ( 
    db *gorm.DB
)

func init() {
    log.Println("Postgres init")
    
    var err error
    db, err = gorm.Open("postgres", "postgres://postgres@db/postgres")
    
    if err != nil {
        log.Fatal(err)
    }
    
    db.DB().SetMaxIdleConns(20)
    db.DB().SetMaxOpenConns(20)
}

func Create(e interface{}) error {
    return db.Create(e).Error
}

type User struct {
	Id    uint
	Name  string
	Email string
}

func Users() ([]User, error) {
	var users []User
	err := db.Order("name asc").Find(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
    
    // return []User{
    //     User { Id: 1, Name: "Name1", Email: "email1@gmail.com"},
    //     User { Id: 2, Name: "Name2", Email: "email2@gmail.com"},
    //     User { Id: 3, Name: "Name3", Email: "email3@gmail.com"},
    //     User { Id: 4, Name: "Name4", Email: "email4@gmail.com"},
    // }, nil
}

