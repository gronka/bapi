package main

import (
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	geohasher "github.com/mmcloughlin/geohash"
	log "github.com/sirupsen/logrus"
)

type ShoutSearchResult struct {
	ShoutUuid gocql.UUID `json:"shoutUuid,omitempty", cql:"shout_uuid"`
}

type SearchShoutsPayload struct {
	Lat     float32             `json:"lat"`
	Lng     float32             `json:"lng"`
	Radius  int                 `json:"radius"`
	Results []ShoutSearchResult `json:"results"`
}

type ShoutPayload struct {
	ShoutUuid    gocql.UUID `json:"shoutUuid,omitempty", cql:"shout_uuid"`
	UserUuid     gocql.UUID `json:"userUuid,omitempty", cql:"user_uuid"`
	ShoutTextNum int        `json:"shoutTextNum", cql:"shout_text_num"`
	Lat          float32    `json:"lat"`
	Lng          float32    `json:"lng"`
	Geohash      string     `json:"geohash,omitempty"`
	Datetime     int64      `json:"datetime"`
}

func hSearchShoutsByNearest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var ssp SearchShoutsPayload
	decodeJson(r, &ssp)

	geohashes := latLngRadiusToGeohashList(ssp.Lat, ssp.Lng, ssp.Radius)
	log.Info(geohashes)
	ssp.searchByGeohashes(geohashes)
	log.Info(ssp.Results)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["results"] = ssp.Results
	gibs.encodeResponse(w)
}

func hCreateShout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var sp ShoutPayload
	decodeJson(r, &sp)

	sp.Geohash = geohasher.Encode(float64(sp.Lat), float64(sp.Lng))
	sp.ShoutUuid, _ = gocql.RandomUUID()
	sp.UserUuid = gibs.UserUuid
	sp.Datetime = time.Now().Unix()
	sp.save()

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}
