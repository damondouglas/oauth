package common

import (
	"log"
	"net/http"
	"runtime/debug"
)

// HandleError writes message and print's stack to log if the error is not nil.
func HandleError(w http.ResponseWriter, err error, message string) {
	if err != nil {
		log.Panic(err)
		debug.PrintStack()
		http.Error(w, message, http.StatusInternalServerError)
	}
}
