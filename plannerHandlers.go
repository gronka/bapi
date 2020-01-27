package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type PlannerRequest struct {
	OwnerUuid         gocql.UUID      `json:"ownerUuid,omitempty", cql:"owner_uuid"`
	DayOfYearStart    int             `json:"dayOfYearStart,omitempty"`
	DayOfYearEnd      int             `json:"dayOfYearEnd,omitempty"`
	AvailableUpdates  int             `json:"availableUpdates", cql:"available_updates"`
	NewAvailableTimes []AvailableTime `json:"newAvailableTimes,omitempty"`
	DeletedTimes      []AvailableTime `json:"deletedTimes,omitempty"`
}

func hGetPlannerByOwnerUuid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl PlannerRequest
	//log.Info("decoding json")
	decodeJson(r, &pl)
	planner := PlannerFromOwnerUuid(pl.OwnerUuid)

	log.Info("planner got from cassandra")

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["planner"] = planner
	gibs.encodeResponse(w)
}

func hRequestFilledTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var fTimeRequest FilledTime
	decodeJson(r, &fTimeRequest)

	//# TODO: tests
	//# reject if too many requests active
	//# -- request filled time with out of range warning
	//# -- request filled time with conflict warning
	//# personal vacation filled time

	// TODO: push notification to requestee
	// TODO: push notification to trainer

	// check permission - currently always returns true
	if !authWriteRequestFilledTime(gibs, fTimeRequest.OwnerUuid) {
		gibs.sendAuthError(w)
		return
	}
	//log.Info(fTimeRequest)
	fTimeRequest.ApptUuid, _ = gocql.RandomUUID()
	fTimeRequest.request()

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hAcceptFilledTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var fTimeRequest FilledTime
	decodeJson(r, &fTimeRequest)

	// TODO: push notification to requestee
	// TODO: push notification to trainer

	// check permission
	if !authPlannerOwner(gibs, fTimeRequest.OwnerUuid) {
		gibs.sendAuthError(w)
		return
	}

	fTimeRequest.accept()
	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hRejectFilledTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var fTimeRequest FilledTime
	decodeJson(r, &fTimeRequest)

	// TODO: push notification to requestee
	// TODO: push notification to trainer

	// check permission
	if !authPlannerOwner(gibs, fTimeRequest.OwnerUuid) {
		gibs.sendAuthError(w)
		return
	}

	fTimeRequest.reject()
	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hCancelFilledTime(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var fTimeRequest FilledTime
	decodeJson(r, &fTimeRequest)

	// TODO: push notification to requestee
	// TODO: push notification to trainer
	fTimeRequest.loadFromKey()

	// check permission
	if !authPlannerOwnerOrRequestee(gibs, fTimeRequest.OwnerUuid, fTimeRequest.RequesteeUuid) {
		gibs.sendAuthError(w)
		return
	}

	if gibs.UserUuid == fTimeRequest.OwnerUuid {
		fTimeRequest.CanceledBy = "owner"
	} else if gibs.UserUuid == fTimeRequest.RequesteeUuid {
		fTimeRequest.CanceledBy = "requestee"
	}

	fTimeRequest.cancel()
	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hGetFilledTimes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: include update counter. If counter does not match client, clear
	// client's cache
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl PlannerRequest
	decodeJson(r, &pl)
	filledTimes := filledTimesFromOwnerUuidTimeRange(pl.OwnerUuid, pl.DayOfYearStart, pl.DayOfYearEnd)

	gibs.ares.Info = ACCEPTED
	// TODO: is this double nested?
	gibs.ares.Body["filledTimes"] = filledTimes
	gibs.encodeResponse(w)
}

func hUpdateAvailableTimes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: make route user-only
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl PlannerRequest
	decodeJson(r, &pl)
	log.Info(pl)

	// check permission
	if !authPlannerOwner(gibs, pl.OwnerUuid) {
		gibs.sendAuthError(w)
		return
	}

	// check that update counters are in sync
	updates := cqlGetPlannerUpdateCounts(pl.OwnerUuid)
	if pl.AvailableUpdates != updates.AvailableUpdates {
		gibs.ares.Info = REJECTED
		gibs.ResponseCode = 500
		gibs.ares.appendFlash("Update error - please refresh.", "error")
		gibs.encodeResponse(w)
		return
	}

	// Increment counter once before committing updates
	cqlIncrementAvailableUpdates(pl.OwnerUuid)

	for _, deletedTime := range pl.DeletedTimes {
		deletedTime.remove(pl.OwnerUuid)
	}

	for _, newAvailableTime := range pl.NewAvailableTimes {
		newAvailableTime.save(pl.OwnerUuid)
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hGetPlannerUpdates(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl PlannerRequest
	decodeJson(r, &pl)
	updates := cqlGetPlannerUpdateCounts(pl.OwnerUuid)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["plannerUpdates"] = updates
	gibs.encodeResponse(w)
}
