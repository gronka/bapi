package main

import (
//"github.com/gocql/gocql"
//log "github.com/sirupsen/logrus"
)

func (pl *EventPayload) dupCreate() {
	stmt := `INSERT INTO user_created_events (user_uuid, event_uuid) VALUES (?, ?)`
	err := s.Query(stmt, pl.UserUuid, pl.EventUuid).Exec()
	check(err, "Failed to insert user_created_events")

	stmt = `INSERT INTO user_admins_events (
	user_uuid,
	end_time,
	start_time,
	event_uuid,
	title,
	tac_name,
	address,
	quick_info) VALUES (?, ?, ?, ?,   ?, ?, ?, ?)`
	err = s.Query(stmt, pl.UserUuid, pl.EndTime, pl.StartTime, pl.EventUuid, pl.Title, pl.TacName, pl.Address, pl.QuickInfo).Exec()
	check(err, "Failed to insert user_admins_events")

	stmt = `INSERT INTO user_rsvped_events (
	user_uuid,
	end_time,
	start_time,
	event_uuid,
	title,
	tac_name,
	address,
	quick_info,
	rsvp) VALUES (?, ?, ?, ?,   ?, ?, ?, ?,   ?)`
	err = s.Query(stmt, pl.UserUuid, pl.EndTime, pl.StartTime, pl.EventUuid, pl.Title, pl.TacName, pl.Address, pl.QuickInfo, "attending").Exec()
	check(err, "Failed to insert user_rsvped_events")

	stmt = `INSERT INTO user_organizing_events (
	user_uuid,
	end_time,
	start_time,
	event_uuid,
	title,
	tac_name,
	address,
	quick_info) VALUES (?, ?, ?, ?,   ?, ?, ?, ?)`
	err = s.Query(stmt, pl.UserUuid, pl.EndTime, pl.StartTime, pl.EventUuid, pl.Title, pl.TacName, pl.Address, pl.QuickInfo).Exec()
	check(err, "Failed to insert user_organizing_events")

	stmt = `INSERT INTO event_has_admins (
	event_uuid,
	user_uuid) VALUES (?, ?)`
	err = s.Query(stmt, pl.EventUuid, pl.UserUuid).Exec()
	check(err, "Failed to insert event_has_admins")
}
