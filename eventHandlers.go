package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//log "github.com/sirupsen/logrus"
)

type EventPayload struct {
	// columns
	EventUuid  gocql.UUID   `json:"eventUuid,omitempty"`
	UserUuid   gocql.UUID   `json:"userUuid,omitempty"`
	TacUuid    gocql.UUID   `json:"tacUuid,omitempty"`
	TacName    string       `json:"tacName,omitempty"`
	Title      string       `json:"title,omitempty"`
	Address    string       `json:"address,omitempty"`
	Lat        float64      `json:"lat,omitempty"`
	Lng        float64      `json:"lng,omitempty"`
	StartTime  int64        `json:"startTime,omitempty"`
	EndTime    int64        `json:"endTime,omitempty"`
	TzOffset   int64        `json:"tzOffset,omitempty"`
	TzId       string       `json:"tzId,omitempty"`
	TzName     string       `json:"tzName,omitempty"`
	LongInfo   string       `json:"longInfo,omitempty"`
	QuickInfo  string       `json:"quickInfo,omitempty"`
	PicUrl     string       `json:"picUrl,omitempty"`
	PinImage   string       `json:"pinImage,omitempty"`
	Admins     []gocql.UUID `json:"admins,omitempty"`
	Organizers []gocql.UUID `json:"organizers,omitempty"`
	// TODO: might replace phone with a json of contact methods
	Phone  string `json:"phone,omitempty"`
	Rating int    `json:"rating,omitempty"`
	Tier   uint64

	// for queries
	OfUser gocql.UUID `json:"ofUser,omitempty"`
}

func hCreateEvent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl EventPayload
	decodeJson(r, &pl)
	// TODO: permission check to create as UserUuid?
	pl.UserUuid = gibs.UserUuid
	pl.create()
	pl.advanceTacLastUsed()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["event"] = pl
	gibs.encodeResponse(w)
}

func hGetUserAdminsEvents(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl EventPayload
	decodeJson(r, &pl)
	// TODO: permission check to see events of user?
	events := pl.getUserAdminsEvents()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["list"] = events
	gibs.encodeResponse(w)
}

func hGetEventByEventUuid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl EventPayload
	decodeJson(r, &pl)
	// TODO: permission check to see events of user?
	event := pl.getEventByEventUuid()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["item"] = event[0]
	gibs.encodeResponse(w)
}
