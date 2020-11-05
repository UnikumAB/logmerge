package utils

import (
	"io"
	"log"
)

func CheckedClose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Fatalf("Cannot close normally: %s", err)
	}
}
