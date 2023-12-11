package handler

import (
	"encoding/json"
	"errors"
	"image-api/pkg/repository"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ListImageMetadata(r http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	data, err := repository.ListImageMetadata(ctx)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Header().Set("Content-Type", "application/json")

	_, err = r.Write(res)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
	}
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
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(r, err.Error(), http.StatusNotFound)
		} else {
			http.Error(r, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	res, err := json.Marshal(data)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}

	r.Header().Set("Content-Type", "application/json")

	_, err = r.Write(res)
	if err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
	}
}
