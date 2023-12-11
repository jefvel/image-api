package handler

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"image"
	"image-api/pkg/repository"
	"image-api/pkg/util"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	metadata, base64Data, err := extractImageData(r)
	if checkError(err, w) {
		return
	}

	result, err := repository.SaveImage(ctx, *metadata, base64Data)
	if checkError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(result)
	if checkError(err, w) {
		return
	}

	_, err = w.Write(res)
	checkError(err, w)
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
	if checkError(err, w) {
		return
	}

	result, err := repository.UpdateImage(ctx, id, *metadata, base64Data)
	if checkError(err, w) {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(result)
	if checkError(err, w) {
		return
	}

	_, err = w.Write(res)
	checkError(err, w)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decodeResponse := r.URL.Query().Has("decode")

	bboxStr := r.URL.Query().Get("bbox")
	var bbox *image.Rectangle
	if len(bboxStr) > 0 {
		bbox, err = util.ParseBBoxString(bboxStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	details, err := repository.GetImageAndMetadata(ctx, id)
	if checkError(err, w) {
		return
	}

	responseData := []byte(*details.Data)
	if bbox != nil {
		imageData, err := base64.StdEncoding.DecodeString(string(*details.Data))
		if checkError(err, w) {
			return
		}

		imgData, err := util.CropImage(imageData, *bbox)
		if checkError(err, w) {
			return
		}

		responseData = []byte(base64.StdEncoding.EncodeToString(imgData))
	}

	if decodeResponse {
		decoded, err := base64.StdEncoding.DecodeString(string(responseData))
		if checkError(err, w) {
			return
		}
		w.Header().Set("Content-Type", "image/"+details.Metadata.Format)
		responseData = decoded
	}

	_, err = w.Write(responseData)

	checkError(err, w)
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

	metadata, err := util.ExtractImageMetadata(imageData)
	if err != nil {
		return nil, "", err
	}

	return metadata, repository.ImageData(base64Data), nil
}

func checkError(err error, w http.ResponseWriter) bool {
	if err == nil {
		return false
	}

	errCode := http.StatusInternalServerError

	if errors.Is(err, repository.ErrNotFound) {
		errCode = http.StatusNotFound
	} else if errors.Is(err, util.ErrInvalidFileFormat) {
		errCode = http.StatusBadRequest
	}

	http.Error(w, err.Error(), errCode)

	return true
}
