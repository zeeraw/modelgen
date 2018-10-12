package db

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

func makeDSN(host, dbname string) string {
	parts := strings.Split(host, "@")
	if len(parts) < 2 {
		log.Fatal(errors.New("invalid connection string format"))
	}

	credentials := strings.Split(parts[0], ":")
	if len(credentials) < 2 {
		log.Fatal(errors.New("invalid connection string format"))
	}
	database := strings.Split(parts[1], ":")
	if len(database) < 2 {
		log.Fatal(errors.New("invalid connection string format"))
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", credentials[0], credentials[1], database[0], database[1], dbname)
}
