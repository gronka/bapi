package main

import (
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type UserSearchResult struct {
	UserUuid       gocql.UUID `json:"userUuid,omitempty", cql:"user_uuid"`
	Fullname       string     `json:"fullname,omitempty"`
	Popularity     float32    `json:"popularity,omitempty"`
	Randomize      float32    `json:"randomize,omitempty"`
	ProfilePicUrl  string     `json:"profilePicUrl,omitempty", cql:"is_videographer"`
	Status         string     `json:"status,omitempty"`
	IsSeller       bool       `json:"isSeller,omitempty", cql:"is_seller"`
	IsTrainer      bool       `json:"isTrainer,omitempty", cql:"is_trainer"`
	IsVideographer bool       `json:"isVideographer,omitempty", cql:"is_videographer"`

	SearchScore int `json:"searchScore"`

	Lat float32 `json:"lat,omitempty"`
	Lng float32 `json:"lng,omitempty"`
}

type SearchUsersPayload struct {
	Lat     float32            `json:"lat"`
	Lng     float32            `json:"lng"`
	Radius  int                `json:"radius"`
	Name    string             `json:"name"`
	Results []UserSearchResult `json:"results"`
}

type UserLocationPayload struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

func hSearchUsersByNearest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var sup SearchUsersPayload
	decodeJson(r, &sup)

	geohashes := latLngRadiusToGeohashList(sup.Lat, sup.Lng, sup.Radius)
	sup.searchByGeohashes(geohashes)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["results"] = sup.Results
	gibs.encodeResponse(w)
}

func hSearchUsersByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var sup SearchUsersPayload
	decodeJson(r, &sup)
	searchName := strings.ToLower(sup.Name)
	words := strings.Fields(searchName)
	if len(words) == 0 {
		gibs.ares.Info = ACCEPTED
		gibs.ares.Body["results"] = sup.Results
		gibs.encodeResponse(w)
		return
	}

	firstWord := words[0]
	// only 1 word entered
	middleWords := firstWord
	lastWord := firstWord
	if len(words) == 2 {
		// Case: 2 words entered, so we don't know if 2nd word is middle or
		// last name
		middleWords = words[1]
		lastWord = words[1]
	} else if len(words) >= 3 {
		middleWords = strings.Join(words[1:(len(words)-1)], " ")
		lastWord = words[len(words)-1]
	}

	// Just keeping as a note: old way
	sup.searchByName("first", firstWord)
	sup.searchByName("middle", middleWords)
	sup.searchByName("last", lastWord)

	// TODO: remove duplicate entries from the results array and assign a
	// higher score to the result (doing this on frontend probably)

	log.Info("user results returned")

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["results"] = sup.Results
	gibs.encodeResponse(w)
}
