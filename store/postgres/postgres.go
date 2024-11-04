package postgres

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"vehicle-plate-recognition/core/environment"
	"vehicle-plate-recognition/store/entities"
)

type Database struct {
	env *environment.Environment
	DB  *gorm.DB
}

func New(env *environment.Environment) *Database {
	var err error
	db, err := gorm.Open(postgres.Open(env.DBURL), &gorm.Config{})
	if err != nil {
		os.Exit(2)
	}
	return &Database{DB: db, env: env}
}

func (db Database) Migrate() {
	err := db.DB.AutoMigrate(
		entities.Vehicle{},
	)
	if err != nil {
		fmt.Println("Failed to migrate tables")
		fmt.Println(err.Error())
		os.Exit(2)
	}
}
