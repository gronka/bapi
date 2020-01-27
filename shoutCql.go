package main

import (
	"github.com/pkg/errors"
)

func (ssp *SearchShoutsPayload) searchByGeohashes(geohashes []string) {
	for _, geohash := range geohashes {
		ssp.searchByGeohash(geohash)
	}
}

func (ssp *SearchShoutsPayload) searchByGeohash(geohash string) {
	query := `SELECT * FROM shout_locs WHERE geohash = ? LIMIT 50`
	iter := s.Query(query, geohash).Iter()

	var shoutRow ShoutSearchResult
	for {
		row := map[string]interface{}{
			"shout_uuid": &shoutRow.ShoutUuid,
		}
		if !iter.MapScan(row) {
			break
		}
		ssp.Results = append(ssp.Results, shoutRow)
	}

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "Error closing iter"))
	}
	return
}

func (sp *ShoutPayload) save() {
	stmt := `INSERT INTO shouts (geohash, shout_uuid, user_uuid, datetime, shout_text_num, lat, lng) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := s.Query(stmt,
		sp.Geohash,
		sp.ShoutUuid,
		sp.UserUuid,
		sp.Datetime,
		sp.ShoutTextNum,
		sp.Lat,
		sp.Lng).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

	sp.saveLoc(sp.Geohash[:3])
	sp.saveLoc(sp.Geohash[:4])
	sp.saveLoc(sp.Geohash[:5])
}

func (sp *ShoutPayload) saveLoc(geohash string) {
	stmt := `INSERT INTO shout_locs (geohash, shout_uuid, lat, lng) VALUES (?, ?, ?, ?)`
	if err := s.Query(stmt, geohash, sp.ShoutUuid, sp.Lat, sp.Lng).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (sp *ShoutPayload) deleet() {
	var oldGeohash string
	stmt := `SELECT geohash FROM shouts WHERE shout_uuid = ?`
	iter := s.Query(stmt, sp.ShoutUuid).Iter()
	iter.Scan(&oldGeohash)

	for i := 3; i <= 5; i++ {
		chars := oldGeohash[:i]
		stmt = `DELETE FROM shout_locs WHERE geohash = ? AND shout_uuid = ? IF EXISTS`
		if err := s.Query(stmt, chars, sp.ShoutUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "error while querying"))
		}
	}
}
