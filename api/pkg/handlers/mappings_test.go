package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	reflect "reflect"
	"testing"
	"urlshortener/pkg/mappings"

	gomock "github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

func TestHandlerRedirect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockMappingRepoInterface(ctrl)
	handler := &MappingHandler{
		MappingRepo: repo,
		Logger:      log.New(os.Stderr, "LOG ", log.Lshortfile),
	}

	mapping := &mappings.Mapping{
		ShortURL: "short_url_test",
		LongURL:  "https://long_url_test",
		Views:    0,
	}

	// Redirect. Correct
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(mapping, nil)
	repo.EXPECT().IncrementMappingViews(mapping).Return(&mappings.Mapping{
		ShortURL: mapping.ShortURL,
		LongURL:  mapping.LongURL,
		Views:    mapping.Views + 1,
	}, nil)

	req := httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w := httptest.NewRecorder()

	handler.Redirect(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("[Redirect. Correct] error")
	}

	// Redirect. No mux.Vars
	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	w = httptest.NewRecorder()

	handler.Redirect(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[Redirect. No mux.Vars] error")
	}

	// Redirect. No URL found
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(nil, mappings.ErrNoURL)

	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w = httptest.NewRecorder()

	handler.Redirect(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[Redirect. No URL found] error")
	}

	// Redirect. Repo GetByShortURL error
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(nil, errors.New("Repo GetByShortURL error"))

	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w = httptest.NewRecorder()

	handler.Redirect(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("[Redirect. Repo error] error")
	}

	// Redirect. Repo IncrementMappingViews error
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(mapping, nil)
	repo.EXPECT().IncrementMappingViews(mapping).Return(nil, errors.New("Repo IncrementMappingViews error"))

	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w = httptest.NewRecorder()

	handler.Redirect(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("[Redirect. IncrementMappingViews error] error")
	}
}

func TestHandlerAdd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockMappingRepoInterface(ctrl)
	handler := &MappingHandler{
		MappingRepo: repo,
		Logger:      log.New(os.Stderr, "LOG ", log.Lshortfile),
	}

	mapping := &mappings.Mapping{
		ShortURL: "short_url_test",
		LongURL:  "https://long_url_test",
		Views:    0,
	}

	mappingBytes, _ := json.Marshal(mapping)

	// Add. Correct
	repo.EXPECT().AddMapping(mapping).Return(mapping, nil)

	req := httptest.NewRequest("POST", "/create", bytes.NewReader(mappingBytes))
	w := httptest.NewRecorder()

	handler.Add(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[Add. Correct] error")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	respMapping := &mappings.Mapping{}
	err := json.Unmarshal(body, respMapping)

	if err != nil || !reflect.DeepEqual(mapping, respMapping) {
		t.Fatalf("[Add. Correct] error")
	}

	// Add. Read request body error
	req = httptest.NewRequest("POST", "/create", &errorReader{})
	w = httptest.NewRecorder()

	handler.Add(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("[Add. Read request body error] error")
	}

	// Add. Invalid JSON
	req = httptest.NewRequest("POST", "/create", bytes.NewReader([]byte("{")))
	w = httptest.NewRecorder()

	handler.Add(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[Add. Invalid JSON] error")
	}

	// Add. Repo AddMapping error
	repo.EXPECT().AddMapping(mapping).Return(nil, errors.New("Repo AddMapping error"))

	req = httptest.NewRequest("POST", "/create", bytes.NewReader(mappingBytes))
	w = httptest.NewRecorder()

	handler.Add(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("[Add. Repo AddMapping error] error")
	}

	// Add. Already Exists error
	repo.EXPECT().AddMapping(mapping).Return(nil, mappings.ErrAlreadyExists)

	req = httptest.NewRequest("POST", "/create", bytes.NewReader(mappingBytes))
	w = httptest.NewRecorder()

	handler.Add(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[Add. Already Exists error] error")
	}

	// Add. Incorrect URL error
	mapping.LongURL = "123"
	mappingBytes, _ = json.Marshal(mapping)

	req = httptest.NewRequest("POST", "/create", bytes.NewReader(mappingBytes))
	w = httptest.NewRecorder()

	handler.Add(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[Add. Incorrect URL error] error")
	}

}

type errorReader struct{}

func (*errorReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("Reader error")
}

func TestHandlerGetMappingInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := NewMockMappingRepoInterface(ctrl)
	handler := &MappingHandler{
		MappingRepo: repo,
		Logger:      log.New(os.Stderr, "LOG ", log.Lshortfile),
	}

	mapping := &mappings.Mapping{
		ShortURL: "short_url_test",
		LongURL:  "https://long_url_test",
		Views:    0,
	}

	// GetMappingInfo. Correct
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(mapping, nil)

	req := httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w := httptest.NewRecorder()

	handler.GetMappingInfo(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("[GetMappingInfo. Correct] error")
	}

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	respMapping := &mappings.Mapping{}
	err := json.Unmarshal(body, respMapping)

	if err != nil || !reflect.DeepEqual(mapping, respMapping) {
		t.Fatalf("[GetMappingInfo. Correct] error")
	}

	// GetMappingInfo. No mux.Vars
	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	w = httptest.NewRecorder()

	handler.GetMappingInfo(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[GetMappingInfo. No mux.Vars] error")
	}

	// GetMappingInfo. No URL found
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(nil, mappings.ErrNoURL)

	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w = httptest.NewRecorder()

	handler.GetMappingInfo(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("[GetMappingInfo. No URL found] error")
	}

	// GetMappingInfo. Repo GetByShortURL error
	repo.EXPECT().GetByShortURL(mapping.ShortURL).Return(nil, errors.New("Repo GetByShortURL error"))

	req = httptest.NewRequest("GET", "/"+mapping.ShortURL, nil)
	req = mux.SetURLVars(req, map[string]string{
		"short_url": mapping.ShortURL,
	})
	w = httptest.NewRecorder()

	handler.GetMappingInfo(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("[Redirect. Repo error] error")
	}
}
