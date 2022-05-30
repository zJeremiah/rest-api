package kittens

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
)

// package kittens is an exmaple api
// you would just delete the kittens folder
// then add the routes that you need for your api

var json = jsoniter.ConfigFastest

type Kitten struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	Fluffy int    `json:"fluffiness"`
	Cute   int    `json:"cuteness"`
}

func Setup() {
	setup.AddEndpoints(
		GetKittensEP(),
		KittenEP(),
	)
}

func GetKittensEP() setup.Endpoint {
	e := setup.Endpoint{
		Name:         "Get All Kittens",
		Version:      "v1",
		Group:        "kittens",
		Path:         "/",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: []Kitten{ // these are example values to be used in the api docs
			{ID: 1, Name: "Fluffums", Breed: "calico", Fluffy: 6, Cute: 7},
			{ID: 2, Name: "Max", Breed: "calico", Fluffy: 5, Cute: 10},
		},
		Description: "This endpoint retrieves all kittens",
		HandlerFunc: GetKittens,
	}

	return e
}

func GetKittens(w http.ResponseWriter, r *http.Request) error {
	var err error
	var respBody []byte

	// Normally you would query a database or api to get the data you needed
	// this is just using the example from the endpoint object
	respBody, err = json.Marshal(GetKittensEP().ResponseBody)
	if err != nil {
		return fmt.Errorf("marshal error for response body %w", err)
	}

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
	return nil
}

func KittenEP() setup.Endpoint {
	e := setup.Endpoint{
		Name:         "Get a Specific Kitten",
		Version:      "v1",
		Group:        "kittens",
		Path:         "/{id}",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: Kitten{ID: 1, Name: "Fluffums", Breed: "calico", Fluffy: 6, Cute: 7}, // example for docs
		Description:  "This endpoint retrieves a specific kitten",
		HandlerFunc:  GetKitten,
		URLParams: []setup.Param{
			{Name: "id", Description: "the id for a kitten"},
		},
	}

	return e
}

func GetKitten(w http.ResponseWriter, r *http.Request) error {
	var err error
	var respBody []byte
	kID := chi.URLParam(r, "id")

	// Normally you would query a database or api to get the data you needed
	// this is just using the example from the endpoint object
	kl := GetKittensEP().ResponseBody.([]Kitten)
	for _, k := range kl {
		if id, _ := strconv.Atoi(kID); k.ID == id {
			respBody, err = json.Marshal(k)
			if err != nil {
				return fmt.Errorf("marshal error for response body %w", err)
			}

			w.Header().Set("Content-Type", setup.ContentJSON)
			w.WriteHeader(http.StatusOK)
			w.Write(respBody)
			return nil
		}
	}

	return &logger.APIErr{
		Internal: logger.Internal{Msg: "kitten id was not found"},
		RespBody: logger.RespBody{Msg: "kitten id was not found: " + kID, Code: 400},
	}
}
