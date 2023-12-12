package handler

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image-api/internal/model"
	"image-api/internal/repository"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	metadata, base64Data, err := extractImageData(r)
	if errorCheck(err, w) {
		return
	}

	result, err := repository.SaveImage(ctx, *metadata, base64Data)
	if errorCheck(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(result)
	if errorCheck(err, w) {
		return
	}

	_, err = w.Write(res)
	errorCheck(err, w)
}

func UpdateImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	metadata, base64Data, err := extractImageData(r)
	if errorCheck(err, w) {
		return
	}

	result, err := repository.UpdateImage(ctx, id, *metadata, base64Data)
	if errorCheck(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(result)
	if errorCheck(err, w) {
		return
	}

	_, err = w.Write(res)
	errorCheck(err, w)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := r.URL.Query()

	decodeResponse := query.Has("decode")

	var bbox *image.Rectangle
	if query.Has("bbox") {
		bboxStr := query.Get("bbox")
		bbox, err = model.ParseBBoxString(bboxStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	details, err := repository.GetImageAndMetadata(ctx, id)
	if errorCheck(err, w) {
		return
	}

	responseData := []byte(*details.Data)
	if bbox != nil {
		imageData, err := base64.StdEncoding.DecodeString(string(*details.Data))
		if errorCheck(err, w) {
			return
		}

		img, err := model.ImageFromBytes(imageData)
		croppedImg, err := img.Crop(*bbox)
		if errorCheck(err, w) {
			return
		}

		responseData = []byte(base64.StdEncoding.EncodeToString(croppedImg.Bytes))
	}

	if decodeResponse {
		decoded, err := base64.StdEncoding.DecodeString(string(responseData))
		if errorCheck(err, w) {
			return
		}

		w.Header().Set("Content-Type", "image/"+details.Metadata.Format)
		responseData = decoded
	}

	_, err = w.Write(responseData)

	errorCheck(err, w)
}

func extractImageData(r *http.Request) (*repository.Metadata, repository.ImageData, error) {
	base64Data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, "", err
	}

	imageData, err := base64.StdEncoding.DecodeString(string(base64Data))
	if err != nil {
		return nil, "", err
	}

	img, err := model.ImageFromBytes(imageData)
	if err != nil {
		return nil, "", err
	}

	metadata := repository.Metadata{
		Width:     img.Width,
		Height:    img.Height,
		Size:      img.Size,
		Format:    img.Format,
		CreatedAt: time.Now(),
	}

	return &metadata, repository.ImageData(base64Data), nil
}
