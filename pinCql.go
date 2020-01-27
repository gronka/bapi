package main

import (
	//"github.com/gocql/gocql"
	log "github.com/sirupsen/logrus"
)

func (pl *EventPayload) pinCreate() {
	msb := MuidSetBuilder{
		Lat:       pl.Lat,
		Lng:       pl.Lng,
		StartTime: pl.StartTime,
		EndTime:   pl.EndTime,
		Tier:      pl.Tier,
		EventUuid: pl.EventUuid.Bytes(),
	}
	log.Info("Generating muids")
	msb.GenerateMuidSet()

	// TODO: lookup if venue has a default pin_image
	stmt := `INSERT INTO event_pins (
		muid,
		event_uuid, 
		end_time, 
		start_time, 
		tac_name, 
		title, 
		lat, 
		lng, 
		tz_offset,
		rating) VALUES (?, ?, ?, ?,  ?, ?, ?, ?,  ?, ?)`

	for _, muid := range msb.Muids {
		log.Trace(muid)
		err := s.Query(stmt,
			muid,
			pl.EventUuid,
			pl.EndTime,
			pl.StartTime,
			pl.TacName,
			pl.Title,
			pl.Lat,
			pl.Lng,
			pl.TzOffset,
			pl.Rating,
		).Exec()
		check(err, "Failed to insert event_pins")
	}

	log.Info("Incrementing counters")
	switch pl.Tier {
	case TIER_BRONZE:
		stmt = `UPDATE muid_counters SET bronze_count = bronze_count + 1 WHERE muid_bin = ?`
	case TIER_SILVER:
		stmt = `UPDATE muid_counters SET silver_count = silver_count + 1 WHERE muid_bin = ?`
	case TIER_GOLD:
		stmt = `UPDATE muid_counters SET gold_count = gold_count + 1 WHERE muid_bin = ?`
	case TIER_DIAMOND:
		stmt = `UPDATE muid_counters SET diamond_count = diamond_count + 1 WHERE muid_bin = ?`
	default:
		panic("tier of event not recognized")
	}

	for _, muidBin := range msb.MuidBins {
		log.Trace(muidBin)
		err := s.Query(stmt, muidBin).Exec()
		check(err, "Failed to update rating_bin counter")
	}

	// need to store these muids for future updates
	stmt = `INSERT INTO event_has_muids (event_uuid, muids, muidBins) VALUES (? ,?, ?)`
	err := s.Query(stmt, pl.EventUuid, msb.Muids, msb.MuidBins).Exec()
	check(err, "Failed to insert event_has_muids")
}
