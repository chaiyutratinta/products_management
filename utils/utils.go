package utils

import (
	"log"
)

//Checker for check error and panic
func Checker(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
