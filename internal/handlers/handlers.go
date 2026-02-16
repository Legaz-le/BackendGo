package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"example.com/mod/internal/models"
	"github.com/go-chi/chi/v5"
)

var bookmarksList = []models.Bookmark{
	{ID: "1", URL: "https://go.dev", Title: "Go Documentation"},
	{ID: "2", URL: "https://github.com", Title: "GitHub"},
}

func findBookmarkIndex(id string) int {
	for i, b := range bookmarksList {
		if b.ID == id {
			return i
		}
	}
	return -1
}

func GetBooks(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query().Get("search")
	check := strings.TrimSpace(params)
	if check == "" {
		json.NewEncoder(w).Encode(bookmarksList)
		return
	} else {
		sliceInfo := []models.Bookmark{}
		checkData := strings.ToLower(check)

		for _, bookmark := range bookmarksList {
			lowerCase := strings.ToLower(bookmark.Title)
			lowerURL := strings.ToLower(bookmark.URL)
			if strings.Contains(lowerCase, checkData) || strings.Contains(lowerURL, checkData) {
				sliceInfo = append(sliceInfo, bookmark)
			}
		}
		json.NewEncoder(w).Encode(sliceInfo)
	}
}

func PostBooks(w http.ResponseWriter, req *http.Request) {
	var newBookmark models.Bookmark
	err := json.NewDecoder(req.Body).Decode(&newBookmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookmarksList = append(bookmarksList, newBookmark)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBookmark)
}

func GetOneBookmark(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	index := findBookmarkIndex(id)
	if index == -1 {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(bookmarksList[index])
}

func UpdateBookmark(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	var updatedBookmark models.Bookmark
	err := json.NewDecoder(req.Body).Decode(&updatedBookmark)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	updatedBookmark.ID = id
	index := findBookmarkIndex(id)
	if index == -1 {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}
	bookmarksList[index] = updatedBookmark
	json.NewEncoder(w).Encode(&updatedBookmark)

}

func DeleteBookmark(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	index := findBookmarkIndex(id)
	if index == -1 {
		http.Error(w, "Bookmark not found", http.StatusNotFound)
		return
	}
	bookmarksList = append(bookmarksList[:index], bookmarksList[index+1:]...)
	w.WriteHeader(http.StatusNoContent)

}