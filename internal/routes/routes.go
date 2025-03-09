package server

import (
	"encoding/json"
	"log"
	"net/http"

	models "github.com/Bharat-kr/url-shortner/internal/models"
	DB "github.com/Bharat-kr/url-shortner/internal/storage"
	"github.com/gorilla/mux"
)

func RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", HelloWorldHandler)
	r.HandleFunc("/shorten", ShortenUrl)
	r.HandleFunc("/urls", GetUrls)

	return r
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

type CreateRequest struct {
	OriginalUrl string `json:"original_url"`
}

func ShortenUrl(w http.ResponseWriter, r *http.Request) {
	var p CreateRequest

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl := GenerateShortCode()
	newUrl := models.Url{OriginalUrl: p.OriginalUrl, ShortUrl: shortUrl}

	// Use the correct database instance to create the new URL record
	if err := DB.DB.Create(&newUrl).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := make(map[string]string)
	resp["original_url"] = p.OriginalUrl
	resp["short_url"] = shortUrl

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func GenerateShortCode() string {

	// TODO: Complete this helper
	code := "jj"

	return code
}

func GetUrls(w http.ResponseWriter, r *http.Request) {
	var urls []models.Url

	if err := DB.DB.Find(&urls).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResp, err := json.Marshal(urls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(jsonResp)
}
