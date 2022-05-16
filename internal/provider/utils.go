package provider

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
)

func debugJson(name string, payload string) {
	var err error
	b := []byte(payload)
	folder := "debug/"

	if _, err := os.Stat(folder); errors.Is(err, os.ErrNotExist) {
		if err = os.Mkdir(folder, os.ModePerm); err != nil {
			log.Panicln("[PANIC] Can't create debug folder ", err)
		}
	}
	if err = os.WriteFile(fmt.Sprintf("%v%v.json", folder, name), b, 0644); err != nil {
		log.Panicln("[PANIC] Can't create debug payload ", err)
	}
}

func loopbackId(loopback string) string {
	re := regexp.MustCompile("[0-9]+")
	data := re.FindAllString(loopback, 1)
	return data[0]
}

func null() []string {
	var null sql.NullString
	empty := []string{}
	empty = append(empty, null.String)
	return empty
}
