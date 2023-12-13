package handler

import (
	"image-api/internal/repository"
	"net/http"
)

func ListImageMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := repository.ListImageMetadata(ctx)
	if errorCheck(err, w) {
		return
	}

	err = writeAsJsonResponse(data, w)
	errorCheck(err, w)
}

func GetImageMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := getIdFromRequest(r)
	if errorCheck(err, w) {
		return
	}

	data, err := repository.GetImageAndMetadata(ctx, id)
	if errorCheck(err, w) {
		return
	}

	err = writeAsJsonResponse(data, w)
	errorCheck(err, w)
}
