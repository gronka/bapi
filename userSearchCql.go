package main

import (
	"github.com/pkg/errors"
)

func (sup *SearchUsersPayload) searchByName(whichName, name string) {
	var query string
	switch whichName {
	case "first":
		query = `SELECT * FROM users WHERE first_name=? LIMIT 30`
	case "middle":
		query = `SELECT * FROM users WHERE middle_name=? LIMIT 30`
	case "last":
		query = `SELECT * FROM users WHERE last_name=? LIMIT 30`
	default:
		// TODO: should this be errors.wrap?
		panic("Incorrect name field")
	}

	if name != "" {
		var userRow UserSearchResult
		iter := s.Query(query, name).Iter()
		for {
			row := map[string]interface{}{
				"user_uuid": &userRow.UserUuid,
			}
			if !iter.MapScan(row) {
				break
			}
			sup.Results = append(sup.Results, userRow)
		}

		if err := iter.Close(); err != nil {
			panic(errors.Wrap(err, "Error closing iter"))
		}
	}

	return
}

func (sup *SearchUsersPayload) searchByGeohashes(geohashes []string) {
	for _, geohash := range geohashes {
		sup.searchByGeohash(geohash)
	}
}

func (sup *SearchUsersPayload) searchByGeohash(geohash string) {
	query := `SELECT * FROM user_locs WHERE geohash = ? LIMIT 50`
	iter := s.Query(query, geohash).Iter()

	var userRow UserSearchResult
	for {
		row := map[string]interface{}{
			"user_uuid": &userRow.UserUuid,
		}
		if !iter.MapScan(row) {
			break
		}
		sup.Results = append(sup.Results, userRow)
	}

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "Error closing iter"))
	}
	return
}
