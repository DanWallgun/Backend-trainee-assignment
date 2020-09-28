package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"urlshortener/pkg/mappings"
	"urlshortener/pkg/utils"
)

// MappingRepoInterface - interface for wrapper
type MappingRepoInterface interface {
	GetByShortURL(string) (*mappings.Mapping, error)
	AddMapping(*mappings.Mapping) (*mappings.Mapping, error)
	IncrementMappingViews(*mappings.Mapping) (*mappings.Mapping, error)
}

// MappingHandler - handles Shorten-Unshorten URL requests
type MappingHandler struct {
	MappingRepo MappingRepoInterface
	Logger      *log.Logger
}

// Redirect - gets short url -> maps to long url -> redirect
func (h *MappingHandler) Redirect(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL, ok := vars["short_url"]
	if !ok {
		ErrorHandler(w, ErrBadShortURL, http.StatusBadRequest)
		return
	}

	mapping, err := h.MappingRepo.GetByShortURL(shortURL)
	if err == mappings.ErrNoURL {
		ErrorHandler(w, ErrURLNotFound, http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Printf("Internal Server Error. %s", err.Error())
		ErrorHandler(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	mapping, err = h.MappingRepo.IncrementMappingViews(mapping)
	if err != nil {
		h.Logger.Printf("Internal Server Error. %s", err.Error())
		ErrorHandler(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, mapping.LongURL, http.StatusSeeOther)
}

// Add - adds a new mapping
func (h *MappingHandler) Add(w http.ResponseWriter, r *http.Request) {
	mapping := new(mappings.Mapping)
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.Logger.Printf("Internal Server Error. %s", err.Error())
		ErrorHandler(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	err = json.Unmarshal(data, mapping)
	if err != nil {
		ErrorHandler(w, ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}
	if !utils.ValidateURL(mapping.LongURL) {
		ErrorHandler(w, ErrIncorrectURL, http.StatusBadRequest)
		return
	}

	mapping, err = h.MappingRepo.AddMapping(mapping)
	if err == mappings.ErrAlreadyExists {
		ErrorHandler(w, ErrAlreadyUsed, http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Printf("Internal Server Error. %s", err.Error())
		ErrorHandler(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(mapping)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// GetMappingInfo retrieves information about given short url: mapping and views
func (h *MappingHandler) GetMappingInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortURL, ok := vars["short_url"]
	if !ok {
		ErrorHandler(w, ErrBadShortURL, http.StatusBadRequest)
		return
	}

	mapping, err := h.MappingRepo.GetByShortURL(shortURL)
	if err == mappings.ErrNoURL {
		ErrorHandler(w, ErrURLNotFound, http.StatusBadRequest)
		return
	}
	if err != nil {
		h.Logger.Printf("Internal Server Error. %s", err.Error())
		ErrorHandler(w, ErrInternalServerError, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(mapping)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
