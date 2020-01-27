package main

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (planner *Planner) loadFromOwnerUuid(ownerUuid gocql.UUID) {
	var availableTimes []AvailableTime
	var availableTime AvailableTime
	query := `SELECT * FROM planner_available_times WHERE owner_uuid=?`
	iter := s.Query(query, ownerUuid).Iter()

	for {
		row := map[string]interface{}{
			"owner_uuid": &availableTime.OwnerUuid,
			"weekday":    &availableTime.Weekday,
			"start_mm":   &availableTime.StartMm,
			"end_mm":     &availableTime.EndMm,
			"tz_offset":  &availableTime.TzOffset,
		}
		if !iter.MapScan(row) {
			break
		}
		availableTimes = append(availableTimes, availableTime)
	}
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}

	// AvailableTimes will be empty if owner has not made modifications
	for _, avTimeFromCassandra := range availableTimes {
		switch avTimeFromCassandra.Weekday {
		case 0:
			planner.Sunday = append(planner.Sunday, avTimeFromCassandra)
		case 1:
			planner.Monday = append(planner.Monday, avTimeFromCassandra)
		case 2:
			planner.Tuesday = append(planner.Tuesday, avTimeFromCassandra)
		case 3:
			planner.Wednesday = append(planner.Wednesday, avTimeFromCassandra)
		case 4:
			planner.Thursday = append(planner.Thursday, avTimeFromCassandra)
		case 5:
			planner.Friday = append(planner.Friday, avTimeFromCassandra)
		case 6:
			planner.Saturday = append(planner.Saturday, avTimeFromCassandra)
		}
	}
}

func (filledTimes *FilledTimes) loadFromOwnerUuid(ownerUuid gocql.UUID) {
	stmt := `SELECT * FROM planner_filled_times WHERE owner_uuid=?`
	iter := s.Query(stmt, ownerUuid).Iter()
	var filledTime FilledTime

	for {
		// TODO: easiest way might be to convert day of the year to date
		row := map[string]interface{}{
			"owner_uuid":     &filledTime.OwnerUuid,
			"day_of_year":    &filledTime.DayOfYear,
			"start_mm":       &filledTime.StartMm,
			"end_mm":         &filledTime.EndMm,
			"reason":         &filledTime.Reason,
			"status":         &filledTime.Status,
			"appt_uuid":      &filledTime.ApptUuid,
			"requestee_uuid": &filledTime.RequesteeUuid,
			"canceled_by":    &filledTime.CanceledBy,
		}
		if !iter.MapScan(row) {
			break
		}
		filledTimes.FilledTimes = append(filledTimes.FilledTimes, filledTime)
	}
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
}

func (filledTime *FilledTime) unpackFromIter(iter *gocql.Iter) {
	for {
		row := map[string]interface{}{
			"owner_uuid":     &filledTime.OwnerUuid,
			"day_of_year":    &filledTime.DayOfYear,
			"start_mm":       &filledTime.StartMm,
			"end_mm":         &filledTime.EndMm,
			"reason":         &filledTime.Reason,
			"status":         &filledTime.Status,
			"tz_offset":      &filledTime.TzOffset,
			"appt_uuid":      &filledTime.ApptUuid,
			"requestee_uuid": &filledTime.RequesteeUuid,
			"canceled_by":    &filledTime.CanceledBy,
		}
		if !iter.MapScan(row) {
			break
		}
	}
}

func (filledTime *FilledTime) loadFromKey() {
	stmt := `SELECT requestee_uuid FROM planner_filled_times WHERE owner_uuid = ? AND day_of_year = ? AND start_mm = ? AND appt_uuid = ?`
	iter := s.Query(stmt, filledTime.OwnerUuid, filledTime.DayOfYear, filledTime.StartMm, filledTime.ApptUuid).Iter()

	filledTime.unpackFromIter(iter)

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
	return
}

func cqlIncrementAvailableUpdates(ownerUuid gocql.UUID) {
	stmt := `UPDATE planner_updates SET available_updates = available_updates + 1 WHERE owner_uuid = ?`
	if err := s.Query(stmt, ownerUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (aTime *AvailableTime) save(ownerUuid gocql.UUID) {
	// Note: pass ownerUuid so that it can't be spoofed. We are running this
	// function over a list, and we are not testing authorization on every loop

	if aTime.Weekday < 0 || aTime.Weekday > 6 ||
		aTime.StartMm < 0 || aTime.StartMm > 1440 ||
		aTime.EndMm < 0 || aTime.EndMm > 1440 ||
		aTime.StartMm > aTime.EndMm {
		log.Error("Incorrect date passed to AvailableTime.save()")
		return
	}

	// Note the use of < for the delete range
	stmt := `DELETE FROM planner_available_times WHERE owner_uuid=? AND weekday=? AND start_mm >= ? AND start_mm < ?`
	if err := s.Query(stmt, ownerUuid, aTime.Weekday, aTime.StartMm, aTime.EndMm).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

	stmt = `INSERT INTO planner_available_times (owner_uuid, weekday, start_mm, end_mm, tz_offset) VALUES (?, ?, ?, ?, ?)`
	if err := s.Query(stmt, ownerUuid, aTime.Weekday, aTime.StartMm, aTime.EndMm, aTime.TzOffset).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (aTime *AvailableTime) remove(ownerUuid gocql.UUID) {
	stmt := `DELETE FROM planner_available_times WHERE owner_uuid=? AND weekday=? AND start_mm=?`
	if err := s.Query(stmt, ownerUuid, aTime.Weekday, aTime.StartMm).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func cqlGetPlannerUpdateCounts(ownerUuid gocql.UUID) (plannerUpdates PlannerUpdates) {
	stmt := `SELECT * FROM planner_updates WHERE owner_uuid=?`
	iter := s.Query(stmt, ownerUuid).Iter()

	iter.Scan(
		&plannerUpdates.OwnerUuid,
		&plannerUpdates.AvailableUpdates,
		&plannerUpdates.FilledUpdates,
	)

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
	return
}

func (fTime *FilledTime) request() {
	if fTime.DayOfYear < 0 || fTime.DayOfYear > 366 ||
		fTime.StartMm < 0 || fTime.StartMm > 1440 ||
		fTime.EndMm < 0 || fTime.EndMm > 1440 ||
		fTime.StartMm > fTime.EndMm {
		log.Error("Incorrect date passed to FilledTime.save()")
		return
	}

	stmt := `INSERT INTO planner_filled_times (owner_uuid, day_of_year, start_mm, end_mm, tz_offset, reason, status, appt_uuid, requestee_uuid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := s.Query(stmt, fTime.OwnerUuid, fTime.DayOfYear, fTime.StartMm, fTime.EndMm, fTime.TzOffset, fTime.Reason, fTime.Status, fTime.ApptUuid, fTime.RequesteeUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

}

func (fTime *FilledTime) accept() {
	status := "accepted"
	stmt := `UPDATE planner_filled_times SET status = ? WHERE owner_uuid = ? AND day_of_year = ? AND start_mm = ? AND appt_uuid = ?`
	if err := s.Query(stmt, status, fTime.OwnerUuid, fTime.DayOfYear, fTime.StartMm, fTime.ApptUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (fTime *FilledTime) reject() {
	status := "rejected"
	stmt := `UPDATE planner_filled_times SET status = ? WHERE owner_uuid = ? AND day_of_year = ? AND start_mm = ? AND appt_uuid = ?`
	if err := s.Query(stmt, status, fTime.OwnerUuid, fTime.DayOfYear, fTime.StartMm, fTime.ApptUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (fTime *FilledTime) cancel() {
	status := "canceled"
	stmt := `UPDATE planner_filled_times SET status = ?, canceled_by = ? WHERE owner_uuid = ? AND day_of_year = ? AND start_mm = ? AND appt_uuid = ?`
	if err := s.Query(stmt, status, fTime.CanceledBy, fTime.OwnerUuid, fTime.DayOfYear, fTime.StartMm, fTime.ApptUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}
