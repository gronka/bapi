package main

import (
	"github.com/pkg/errors"
)

func (scp *SearchCoursesPayload) searchByGeohashes(geohashes []string) {
	for _, geohash := range geohashes {
		scp.searchByGeohash(geohash)
	}
}

func (scp *SearchCoursesPayload) searchByGeohash(geohash string) {
	query := `SELECT * FROM course_locs WHERE geohash = ? LIMIT 50`
	iter := s.Query(query, geohash).Iter()

	var courseRow CourseSearchResult
	for {
		row := map[string]interface{}{
			"course_uuid": &courseRow.CourseUuid,
		}
		if !iter.MapScan(row) {
			break
		}
		scp.Results = append(scp.Results, courseRow)
	}

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "Error closing iter"))
	}
	return
}
