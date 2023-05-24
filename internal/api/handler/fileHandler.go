package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	blurhash "github.com/buckket/go-blurhash"
	"github.com/gorilla/mux"
	"github.com/souravlayek/file-hosting/internal/database"
	"github.com/souravlayek/file-hosting/internal/model"
)

type UploadResponse struct {
	Url      string `json:"url"`
	BlurHash string `json:"blurhash"`
}

var allowedMimeTypes = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"application/pdf": true,
}

// FileHandler handles the file upload
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	mimeTypes := handler.Header["Content-Type"]
	if !allowedMimeTypes[mimeTypes[0]] {
		fmt.Println("Invalid file type")
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}
	// store in media folder
	fileName := handler.Filename
	extension := strings.Split(fileName, ".")[1]
	fileNameWithoutExtension := strings.Split(fileName, ".")[0]
	fileNameToStore := fileNameWithoutExtension + "-" + time.Now().Format("_2006_01_02_15_04_05") + "." + extension
	os.MkdirAll("media/", os.ModePerm)
	f, err := os.OpenFile("media/"+fileNameToStore, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var _blurhash string
	if mimeTypes[0] == "image/jpeg" || mimeTypes[0] == "image/png" {
		// generate blurhash
		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		img, _, err := image.Decode(file)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_blurhash, err = blurhash.Encode(4, 4, img)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	myMetaData := model.MetaData{
		ID:        primitive.NewObjectID().Hex(),
		FileName:  fileNameToStore,
		MimeType:  mimeTypes[0],
		Size:      handler.Size,
		BlurHash:  _blurhash,
		FilePath:  "media/" + fileNameToStore,
		CreatedAt: time.Now().Unix(),
	}

	_, err = database.DB.Collection("metaData").InsertOne(context.TODO(), myMetaData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hostname := os.Getenv("ENDPOINT")
	if hostname == "" {
		http.Error(w, "Hostname not set", http.StatusInternalServerError)
		return
	}
	// return the url and blurhash
	response := UploadResponse{
		Url:      hostname + "s/" + myMetaData.ID,
		BlurHash: _blurhash,
	}
	json.NewEncoder(w).Encode(response)
}

// FileDownloadHandler handles the file download
func FileDownloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var myMetaData model.MetaData
	err := database.DB.Collection("metaData").FindOne(context.TODO(), bson.M{
		"_id": id,
	}).Decode(&myMetaData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// open file
	f, err := os.Open(myMetaData.FilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer f.Close()
	// serve file
	http.ServeFile(w, r, myMetaData.FilePath)
}
