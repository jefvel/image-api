package handler

import (
	"errors"
	"image-api/internal/model"
	"image-api/internal/repository"
	"net/http"
)

func errorCheck(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	errCode := http.StatusInternalServerError

	if errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
	} else if errors.Is(err, model.ErrInvalidFileFormat) {
		errCode = http.StatusBadRequest
	}

	http.Error(w, err.Error(), errCode)

	return true
}
