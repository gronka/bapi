package main

import (
	"github.com/gocql/gocql"
	//log "github.com/sirupsen/logrus"
)

func (gibs *Gibs) isEventAdmin(eventUuid gocql.UUID) bool {
	var gotUserUuid gocql.UUID
	query := `SELECT user_uuid FROM event_has_admins WHERE event_uuid=? AND user_uuid=?`
	iter := s.Query(query, eventUuid, gibs.UserUuid).Iter()
	iter.Scan(&gotUserUuid)

	if gotUserUuid != EmptyUuidBytes {
		return true
	}
	return false
}
