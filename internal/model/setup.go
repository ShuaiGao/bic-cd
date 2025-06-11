package model

import "bic-cd/pkg/db"

func Setup() {
	err := db.DB().AutoMigrate(&User{}, &Service{}, &ServiceInstance{})
	if err != nil {
		panic(err)
	}
}
