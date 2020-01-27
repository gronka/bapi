package main

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
)

// If any changes are made to a planner, increment this
type PlannerUpdates struct {
	OwnerUuid        gocql.UUID `json:"ownerUuid,omitempty", cql:"owner_uuid"`
	AvailableUpdates int        `json:"availableUpdates", cql:"available_updates"`
	FilledUpdates    int        `json:"filledUpdates", cql:"filled_updates"`
}

type Planner struct {
	OwnerUuid gocql.UUID      `json:"ownerUuid,omitempty", cql:"owner_uuid"`
	Monday    []AvailableTime `json:"monday"`
	Tuesday   []AvailableTime `json:"tuesday"`
	Wednesday []AvailableTime `json:"wednesday"`
	Thursday  []AvailableTime `json:"thursday"`
	Friday    []AvailableTime `json:"friday"`
	Saturday  []AvailableTime `json:"saturday"`
	Sunday    []AvailableTime `json:"sunday"`
	//FilledTimes []FilledTime    `json:"filledTime"`
}

type AvailableTime struct {
	OwnerUuid gocql.UUID `json:"ownerUuid,omitempty", cql:"owner_uuid"`
	Weekday   int        `json:"weekday"`
	StartMm   int        `json:"startMm", cql:"start_mm"`
	EndMm     int        `json:"endMm", cql:"end_mm"`
	TzOffset  int        `json:"tzOffset", cql:"tz_offset"`
}

type FilledTimes struct {
	FilledTimes []FilledTime `json:"filledTimes"`
}

type FilledTime struct {
	OwnerUuid     gocql.UUID `json:"ownerUuid,omitempty", cql:"owner_uuid"`
	DayOfYear     int        `json:"dayOfYear", cql:"day_of_year"`
	StartMm       int        `json:"startMm", cql:"start_mm"`
	EndMm         int        `json:"endMm", cql:"end_mm"`
	TzOffset      int        `json:"tzOffset", cql:"tz_offset"`
	Reason        string     `json:"reason"`
	Status        string     `json:"status"`
	ApptUuid      gocql.UUID `json:"apptUuid,omitempty", cql:"appt_uuid"`
	RequesteeUuid gocql.UUID `json:"requesteeUuid,omitempty", cql:"requestee_uuid"`
	CanceledBy    string     `json:"canceledBy,omitempty", cql:"canceled_by"`
}

// TODO: pull filledTimes by current week in view/daterange

func PlannerFromOwnerUuid(ownerUuid gocql.UUID) (planner Planner) {
	planner.loadFromOwnerUuid(ownerUuid)
	return
}

func filledTimesFromOwnerUuidTimeRange(ownerUuid gocql.UUID, dayOfYearStart, dayOfYearEnd int) (filledTimes FilledTimes) {
	var filledTime FilledTime

	dayOfYearStart--

	stmt := `SELECT * FROM planner_filled_times WHERE owner_uuid=? and day_of_year >= ? AND day_of_year <= ?`

	iter := s.Query(stmt, ownerUuid, dayOfYearStart, dayOfYearEnd).Iter()

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
		filledTimes.FilledTimes = append(filledTimes.FilledTimes, filledTime)

	}
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
	return
}

func requestFilledTimes(requesteeUuid, trainerUuid gocql.UUID, weekday, startTime, endTime int) {
	status := REQUEST_PENDING
	apptUuid, _ := gocql.RandomUUID()
	stmt := `INSERT INTO planner_filled_times (owner_uuid, weekday, startTime, endTime, status, appt_uuid, player_uuid) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := s.Query(stmt,
		requesteeUuid,
		weekday,
		startTime,
		endTime,
		status,
		apptUuid,
		trainerUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

}
