package handler

import (
	"encoding/json"
	"errors"
	"image-api/internal/model"
	"image-api/internal/repository"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var (
	ErrInvalidID = errors.New("Invalid ID")
)

func getIdFromRequest(r *http.Request) (int, error) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		return 0, ErrInvalidID
	}

	return id, nil
}

func writeAsJsonResponse(data any, w http.ResponseWriter) error {
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)

	return err
}

// errorCheck handles error response writing for HTTP requests.
// It logs the error and writes an appropriate HTTP error response.
func errorCheck(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	var errCode int
	errMsg := err.Error()

	switch {
	case errors.Is(err, repository.ErrNotFound):
		errCode = http.StatusNotFound
	case errors.Is(err, model.ErrInvalidFileFormat):
		errCode = http.StatusBadRequest
	case errors.Is(err, ErrInvalidID):
		errCode = http.StatusBadRequest
	default:
		errMsg = "Internal server error"
		errCode = http.StatusInternalServerError
	}

	log.Printf("Error: %v", err)
	http.Error(w, errMsg, errCode)

	return true
}
