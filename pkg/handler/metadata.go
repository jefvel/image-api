package handler

import (
	"encoding/json"
	"image-api/pkg/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ListImageMetadata(r http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	data, err := repository.ListImageMetadata(ctx)
	if checkError(err, r) {
		return
	}

	res, err := json.Marshal(data)
	if checkError(err, r) {
		return
	}

	r.Header().Set("Content-Type", "application/json")

	_, err = r.Write(res)
	checkError(err, r)
}

func GetImageMetadata(r http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	params := mux.Vars(req)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(r, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := repository.GetImageMetadata(ctx, id)
	if checkError(err, r) {
		return
	}

	res, err := json.Marshal(data)
	if checkError(err, r) {
		return
	}

	r.Header().Set("Content-Type", "application/json")
	_, err = r.Write(res)

	checkError(err, r)
}
