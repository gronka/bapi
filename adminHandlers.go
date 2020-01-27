package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

func hTruncateAllTables(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: make route admin
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)

	for _, table := range TABLE_NAMES {
		stmt := `TRUNCATE ` + table
		if err := s.Query(stmt).Exec(); err != nil {
			panic(errors.Wrap(err, "error while truncating"))
		}
	}

	stmt := `INSERT INTO kapi.users (user_uuid) VALUES (99999999-9999-9999-9999-999999999999);`
	if err := s.Query(stmt).Exec(); err != nil {
		panic(errors.Wrap(err, "error while adding user9"))
	}
	stmt = `UPDATE kapi.users SET 
			phone_verified = true,
			calling_code = '1',
			iso_phone = 'us9',
			naked_phone = '9',
			phone_country_iso = 'us',
			account_status = 'FINE',
			fullname = 'Terry Hacker',
			about = 'I made my account second',
			rating = 3.4,
			password = '9'
		WHERE user_uuid=99999999-9999-9999-9999-999999999999;
		`
	if err := s.Query(stmt).Exec(); err != nil {
		panic(errors.Wrap(err, "error while adding user9"))
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}
