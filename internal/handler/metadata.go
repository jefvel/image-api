package handler

import (
	"encoding/json"
	"image-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ListImageMetadata(r http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	data, err := repository.ListImageMetadata(ctx)
	if errorCheck(err, r) {
		return
	}

	res, err := json.Marshal(data)
	if errorCheck(err, r) {
		return
	}

	r.Header().Set("Content-Type", "application/json")

	_, err = r.Write(res)
	errorCheck(err, r)
}

func GetImageMetadata(r http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := repository.GetImageAndMetadata(ctx, id)
	if errorCheck(err, r) {
		return
	}

	res, err := json.Marshal(data)
	if errorCheck(err, r) {
		return
	}

	r.Header().Set("Content-Type", "application/json")
	_, err = r.Write(res)

	errorCheck(err, r)
}
