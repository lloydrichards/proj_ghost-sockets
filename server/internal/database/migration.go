package database

import (
	"log"
)

func (s *service) migrate() error {
	log.Println("Running migrations...")

	s.db.AutoMigrate(&User{})
	s.db.AutoMigrate(&Session{})

	return nil
}
