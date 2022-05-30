package kittns

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/rest-api/internal/logger"
	"github.com/rest-api/internal/setup"
)

// package kittns is an exmaple api
// you would just delete the kittns folder
// then add the routes that you need for your api

var json = jsoniter.ConfigFastest

type Kitten struct {
	ID     int    `json:"id,omitempty"`
	Name   string `json:"name"`
	Breed  string `json:"breed"`
	Fluffy int    `json:"fluffiness"`
	Cute   int    `json:"cuteness"`
}

func Setup() {
	setup.AddEndpoints(
		GetKittensEP(),
		KittenEP(),
		RMKittenEP(),
		AddKittnEP(),
	)
}

func GetKittensEP() setup.Endpoint {
	e := setup.Endpoint{
		Name:         "Get All Kittens",
		Version:      "v1",
		Group:        "kittns",
		Path:         "/",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: []Kitten{ // these are example values to be used in the api docs
			{ID: 1, Name: "Fluffums", Breed: "calico", Fluffy: 6, Cute: 7},
			{ID: 2, Name: "Max", Breed: "calico", Fluffy: 5, Cute: 10},
		},
		Description: "This endpoint retrieves all kittns",
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
		Group:        "kittns",
		Path:         "/{id}",
		Method:       setup.GET,
		ResponseType: setup.ContentJSON,
		ResponseBody: Kitten{ID: 1, Name: "Fluffums", Breed: "calico", Fluffy: 6, Cute: 7}, // example for docs
		Description:  "This endpoint retrieves a specific kittn",
		HandlerFunc:  GetKitten,
		URLParams: []setup.Param{
			{Name: "id", Description: "the id for a kittn"},
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
		Internal: logger.Internal{Msg: "kittn id was not found"},
		RespBody: logger.RespBody{Msg: "kittn id was not found: " + kID, Code: 400},
	}
}

func RMKittenEP() setup.Endpoint {
	e := setup.Endpoint{
		Name:         "Delete a Specific Kitten",
		Version:      "v1",
		Group:        "kittns",
		Path:         "/{id}",
		Method:       setup.DELETE,
		ResponseType: setup.ContentJSON,
		Description:  "This endpoint deletes a specific kittn",
		HandlerFunc:  RMKitten,
		URLParams: []setup.Param{
			{Name: "id", Description: "the id for a kittn"},
		},
	}

	return e
}

func RMKitten(w http.ResponseWriter, r *http.Request) error {
	// Normally you would send a delete request to a database or api
	// this actually doesn't do anything but lie to you :)

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func AddKittnEP() setup.Endpoint {
	e := setup.Endpoint{
		Name:        "Add a New Kittn",
		Version:     "v1",
		Group:       "kittns",
		Path:        "/",
		Method:      setup.POST,
		RequestType: setup.ContentJSON,
		RequestBody: Kitten{
			Name:   "Stealth",
			Breed:  "Siamese",
			Fluffy: 2,
			Cute:   3,
		},
		ResponseBody: Kitten{
			ID:     3,
			Name:   "Stealth",
			Breed:  "Siamese",
			Fluffy: 2,
			Cute:   3,
		},
		ResponseType: setup.ContentJSON,
		Description:  "This endpoint deletes a specific kittn",
		HandlerFunc:  AddKittn,
	}

	return e
}

func AddKittn(w http.ResponseWriter, r *http.Request) error {

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		return logger.NewError(err.Error(), "could not read request body", 406, err)
	}
	k := Kitten{}
	err = json.Unmarshal(buf, &k)
	if err != nil {
		return logger.NewError(err.Error(), "could not parse request body", 406, err)
	}

	// Normally you would send a insert request to a database or
	// a create request to another api, this one doesn't do anything
	k.ID = 3
	respBody, err := json.Marshal(k)
	if err != nil {
		return logger.NewError("there has been a problem", "could not marshal response body", 500, err)
	}

	w.Header().Set("Content-Type", setup.ContentJSON)
	w.WriteHeader(http.StatusAccepted)
	w.Write(respBody)
	return nil
}
