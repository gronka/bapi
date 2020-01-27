package main

import (
	"math/rand"

	"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

func (pl *EventPayload) advanceTacLastUsed() {
	now := timeNowMilli()
	stmt := `UPDATE user_has_tacs SET last_used=? WHERE user_uuid=? AND tac_uuid=?`
	err := s.Query(stmt, now, pl.UserUuid, pl.TacUuid).Exec()
	check(err, "failed to update last_used of user_has_tacs")
}

func (pl *EventPayload) create() {
	// TODO: get tier from hosting/org account
	pl.Tier = TIER_BRONZE
	// TODO: calculate rating from amount paid, tier, and maybe popularity
	pl.Rating = rand.Intn(100)

	pl.EventUuid, _ = gocql.RandomUUID()
	stmt := `INSERT INTO events (
		event_uuid, 
		tac_uuid, 
		tac_name, 
		title, 
		address, 
		lat, 
		lng, 
		tz_offset, 
		tz_id,
		tz_name,
		start_time, 
		end_time, 
		quick_info, 
		admins, 
		phone,
		tier,
		rating) VALUES (?, ?, ?, ?,   ?, ?, ?, ?,   ?, ?, ?, ?,   ?, ?, ?, ?,   ?)`

	err := s.Query(stmt,
		pl.EventUuid,
		pl.TacUuid,
		pl.TacName,
		pl.Title,
		pl.Address,
		pl.Lat,
		pl.Lng,
		pl.TzOffset,
		pl.TzId,
		pl.TzName,
		pl.StartTime,
		pl.EndTime,
		pl.QuickInfo,
		[]gocql.UUID{pl.UserUuid},
		pl.Phone,
		pl.Tier,
		pl.Rating,
	).Exec()
	check(err, "Failed to insert into events")

	log.Info("Creating dups")
	pl.dupCreate()
	log.Info("Creating pins")
	pl.pinCreate()
}

func (pl *EventPayload) getUserAdminsEvents() []EventPayload {
	query := `SELECT * FROM user_admins_events WHERE user_uuid=?`
	iter := s.Query(query, pl.OfUser).Iter()

	var events []EventPayload
	var buf EventPayload
	for {
		row := map[string]interface{}{
			"event_uuid": &buf.EventUuid,
			"tac_name":   &buf.TacName,
			"title":      &buf.Title,
			"address":    &buf.Address,
			"start_time": &buf.StartTime,
			"end_time":   &buf.EndTime,
			"quick_info": &buf.QuickInfo,
		}
		if !iter.MapScan(row) {
			break
		}
		events = append(events, buf)
	}
	err := iter.Close()
	checkIter(err)
	return events
}

func (pl *EventPayload) getEventsAttendingByUserUuid() []EventPayload {
	query := `SELECT * FROM events WHERE user_uuid=?`
	iter := s.Query(query, pl.OfUser).Iter()

	var events []EventPayload
	var buf EventPayload
	for {
		row := map[string]interface{}{
			"event_uuid": &buf.EventUuid,
			"tac_name":   &buf.TacName,
			"title":      &buf.Title,
			"address":    &buf.Address,
			"start_time": &buf.StartTime,
			"end_time":   &buf.EndTime,
			"quick_info": &buf.QuickInfo,
			"admins":     &buf.Admins,
		}
		if !iter.MapScan(row) {
			break
		}
		events = append(events, buf)
	}
	err := iter.Close()
	checkIter(err)
	return events
}

func (pl *EventPayload) getEventByEventUuid() []EventPayload {
	query := `SELECT * FROM events WHERE event_uuid=?`
	iter := s.Query(query, pl.EventUuid).Iter()

	log.Info(pl.EventUuid)
	log.Info(pl)

	var events []EventPayload
	var buf EventPayload
	for {
		row := map[string]interface{}{
			"event_uuid": &buf.EventUuid,
			"tac_uuid":   &buf.TacUuid,
			"tac_name":   &buf.TacName,
			"lat":        &buf.Lat,
			"lng":        &buf.Lng,
			"title":      &buf.Title,
			"address":    &buf.Address,
			"start_time": &buf.StartTime,
			"end_time":   &buf.EndTime,
			"quick_info": &buf.QuickInfo,
			"admins":     &buf.Admins,
		}
		if !iter.MapScan(row) {
			break
		}
		events = append(events, buf)
	}
	err := iter.Close()
	checkIter(err)
	return events
}

func (pl *EventPayload) changeLocation(tac UserHasTacsModel) {
	stmt := `UPDATE events SET 
		tac_uuid=?,
		tac_name=?,
		address=?,
		lat=?,
		lng=?,
		tz_offset=?,
		tz_id=? WHERE event_uuid=?`
	err := s.Query(stmt, tac.TacUuid, tac.Name, tac.Address, tac.Lat, tac.Lng, tac.TzOffset, tac.TzId, pl.EventUuid).Exec()
	check(err, "failed to update events with new location")
}
