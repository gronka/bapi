package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type EventCapsule struct {
	EventUuid gocql.UUID  `json:"eventUuid"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`

	// internals
	UserUuid gocql.UUID
	EndTime  int64
}

func (ecap *EventCapsule) getEndTime() {
	query := `SELECT end_time FROM events WHERE event_uuid=?`
	iter := s.Query(query, ecap.EventUuid).Iter()
	iter.Scan(&ecap.EndTime)
}

func (ecap *EventCapsule) changeLongInfo(longInfo string) {
	stmt := `UPDATE events SET long_info=? WHERE event_uuid=?`
	err := s.Query(stmt, longInfo, ecap.EventUuid).Exec()
	check(err, "failed to update long_info in events")
}

func (ecap *EventCapsule) changeTitle(title string) {
	ecap.getEndTime()

	stmt := `UPDATE events SET title=? WHERE event_uuid=?`
	err := s.Query(stmt, title, ecap.EventUuid).Exec()
	check(err, "failed to update title in events")

	stmt = `UPDATE user_admins_events SET title=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, title, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update title in user_admins_events")

	stmt = `UPDATE user_rsvped_events SET title=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, title, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update title in user_rsvped_events")

	stmt = `UPDATE user_organizing_events SET title=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, title, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update title in user_organizing_events")

	// TODO: update all pins
}

func (ecap *EventCapsule) changeQuickInfo(quickInfo string) {
	ecap.getEndTime()

	stmt := `UPDATE events SET quick_info=? WHERE event_uuid=?`
	err := s.Query(stmt, quickInfo, ecap.EventUuid).Exec()
	check(err, "failed to update quick_info in events")

	stmt = `UPDATE user_admins_events SET quick_info=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, quickInfo, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update quick_info in user_admins_events")

	stmt = `UPDATE user_rsvped_events SET quick_info=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, quickInfo, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update quick_info in user_rsvped_events")

	stmt = `UPDATE user_organizing_events SET quick_info=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, quickInfo, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update quick_info in user_organizing_events")

	// TODO: update all pins
}

func (ecap *EventCapsule) changeStartTime(startTime int64) {
	ecap.getEndTime()

	stmt := `UPDATE events SET start_time=? WHERE event_uuid=?`
	err := s.Query(stmt, startTime, ecap.EventUuid).Exec()
	check(err, "failed to update start_time in events")

	stmt = `UPDATE user_admins_events SET start_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, startTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update start_time in user_admins_events")

	stmt = `UPDATE user_rsvped_events SET start_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, startTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update start_time in user_rsvped_events")

	stmt = `UPDATE user_organizing_events SET start_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, startTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update start_time in user_organizing_events")

	// TODO: update all pins
}

func (ecap *EventCapsule) changeEndTime(endTime int64) {
	ecap.getEndTime()

	stmt := `UPDATE events SET end_time=? WHERE event_uuid=?`
	err := s.Query(stmt, endTime, ecap.EventUuid).Exec()
	check(err, "failed to update end_time in events")

	stmt = `UPDATE user_admins_events SET end_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, endTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update end_time in user_admins_events")

	stmt = `UPDATE user_rsvped_events SET end_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, endTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update end_time in user_rsvped_events")

	stmt = `UPDATE user_organizing_events SET end_time=? WHERE event_uuid=? AND end_time=? AND user_uuid=?`
	err = s.Query(stmt, endTime, ecap.EventUuid, ecap.EndTime, ecap.UserUuid).Exec()
	check(err, "failed to update end_time in user_organizing_events")

	// TODO: update all pins
}

func hEventFieldUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var ecap EventCapsule
	decodeJson(r, &ecap)

	if authorized := gibs.authUserCanAdminEvent(ecap.EventUuid); !authorized {
		gibs.sendAuthError(w)
		return
	}
	ecap.UserUuid = gibs.UserUuid

	var str, failMsg string
	var timeValue int64

	switch ecap.Field {
	case "title":
		str, failMsg = SanitizeString(ecap.Value, 3, 80)
	case "quickInfo":
		str, failMsg = SanitizeString(ecap.Value, 3, 140)
	case "longInfo":
		str, failMsg = SanitizeString(ecap.Value, 3, 2000)
	case "startTime":
		timeValue, failMsg = SanitizeTime(ecap.Value)
	case "endTime":
		timeValue, failMsg = SanitizeTime(ecap.Value)
	default:
		failMsg = "Could not update field with name: " + ecap.Field
	}

	if failMsg != "" {
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash(failMsg, "error")
		return
	}

	switch ecap.Field {
	case "title":
		ecap.changeTitle(str)
	case "quickInfo":
		ecap.changeQuickInfo(str)
	case "longInfo":
		ecap.changeLongInfo(str)
	case "startTime":
		ecap.changeStartTime(timeValue)
	case "endTime":
		ecap.changeEndTime(timeValue)
	default:
		panic("impossible case reached in hEventFieldUpdate")
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)

}

func hEventLocationChange(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl EventPayload
	decodeJson(r, &pl)
	if authorized := gibs.authUserCanAdminEvent(pl.EventUuid); !authorized {
		log.Info("not authorized")
		gibs.sendAuthError(w)
		return
	}
	var tac UserHasTacsModel
	tac.UserUuid = gibs.UserUuid
	tac.TacUuid = pl.TacUuid
	tac.load()
	log.Info("tac")
	log.Info(tac)
	pl.changeLocation(tac)

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}
