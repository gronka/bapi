package main

import (
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (user *User) unpackFromIter(iter *gocql.Iter) {
	for {
		row := map[string]interface{}{
			"user_uuid":         &user.UserUuid,
			"email":             &user.Email,
			"naked_phone":       &user.NakedPhone,
			"phone_country_iso": &user.PhoneCountryIso,
			"iso_phone":         &user.IsoPhone,
			"calling_code":      &user.CallingCode,
			"password":          &user.Password,
			"fullname":          &user.Fullname,
			"phone_verified":    &user.PhoneVerified,
			"email_verified":    &user.EmailVerified,
			"account_status":    &user.AccountStatus,
		}
		if !iter.MapScan(row) {
			break
		}
	}
}

func (user *User) loadFromUserUuid(userUuid gocql.UUID) {
	var iter *gocql.Iter

	stmt := `SELECT user_uuid, email, naked_phone, phone_country_iso, iso_phone, calling_code, password, fullname, account_status FROM users WHERE user_uuid=?`
	iter = s.Query(stmt, userUuid).Iter()

	user.unpackFromIter(iter)

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
	return
}

func (user *User) loadFromSignInForm(authType, isoPhone, email string) {
	var stmt, value string

	if authType == AUTH_IS_PHONE {
		stmt = `SELECT user_uuid, email, naked_phone, phone_country_iso, iso_phone, calling_code, password, fullname, phone_verified, email_verified, account_status FROM users WHERE iso_phone=?`
		value = isoPhone

	} else if authType == AUTH_IS_EMAIL {
		stmt = `SELECT user_uuid, email, naked_phone, phone_country_iso, iso_phone, calling_code, password, fullname, account_status FROM users WHERE email=?`
		value = email
	}

	iter := s.Query(stmt, value).Iter()
	user.unpackFromIter(iter)

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}

	user.AuthType = authType
	return
}

func (user *User) saveOnSignUp() {
	stmt := `INSERT INTO users (user_uuid, email, naked_phone, phone_country_iso, iso_phone, calling_code, password, phone_verified, email_verified, account_status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	if err := s.Query(stmt,
		user.UserUuid,
		user.Email,
		user.NakedPhone,
		user.PhoneCountryIso,
		user.IsoPhone,
		user.CallingCode,
		user.Password,
		user.PhoneVerified,
		user.EmailVerified,
		user.AccountStatus).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (user *User) lookupAccountStatus() string {
	var iter *gocql.Iter
	var lookup User
	lookup.AccountStatus = DOES_NOT_EXIST

	// TODO: add uuid?
	if user.AuthType == AUTH_IS_PHONE {
		iter = s.Query(`SELECT "account_status" FROM users
		WHERE iso_phone=?`, user.IsoPhone).Iter()
		iter.Scan(
			&lookup.AccountStatus,
		)
	} else if user.AuthType == AUTH_IS_EMAIL {
		iter = s.Query(`SELECT account_status FROM users
		WHERE email=?`, user.Email).Iter()
		iter.Scan(
			&lookup.AccountStatus,
		)
	}

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}

	return lookup.AccountStatus
}

func getAllUsers() []interface{} {
	query := `SELECT * FROM users`
	iter := s.Query(query).Iter()

	var users []interface{}
	res := make(map[string]interface{})
	for iter.MapScan(res) {
		users = append(users, res)
		res = make(map[string]interface{})
	}

	//var user User
	//var fullname string
	//for iter.Scan(&fullname) {
	//users = append(users, fullname)
	//}

	// results fail to marshal to json with this method - uses cql column names
	//results, _ := iter.SliceMap()

	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "Error closing iter"))
	}
	return users
}

func (user *User) saveLoc(geohash string, lat, lng float32) {
	stmt := `INSERT INTO user_locs (geohash, user_uuid, lat, lng) VALUES (?, ?, ?, ?)`
	if err := s.Query(stmt, geohash, user.UserUuid, lat, lng).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (user *User) removeLocs() {
	var oldGeohash string
	stmt := `SELECT geohash FROM users WHERE user_uuid = ?`
	iter := s.Query(stmt, user.UserUuid).Iter()
	iter.Scan(&oldGeohash)

	// TODO: batch statements? pass value from SELECT to DELETE
	log.Info(oldGeohash)

	if oldGeohash != "" {
		log.Info(oldGeohash)
		for i := 3; i <= 5; i++ {
			chars := oldGeohash[:i]
			stmt = `DELETE FROM user_locs WHERE geohash = ? AND user_uuid = ? IF EXISTS`
			if err := s.Query(stmt, chars, user.UserUuid).Exec(); err != nil {
				panic(errors.Wrap(err, "error while querying"))
			}

		}
	}
}

func (user *User) saveLatLngGeo(lat, lng float32, geohash string) {
	log.Info("++++++saving geoash")
	log.Info(geohash)
	stmt := `UPDATE users SET lat = ?, lng = ?, geohash = ? WHERE user_uuid = ?`
	if err := s.Query(stmt, lat, lng, geohash, user.UserUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (user *User) removeByUserUuid() {
	stmt := `DELETE FROM users WHERE user_uuid=?`
	if err := s.Query(stmt, user.UserUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}
}

func (user *User) removeByIsoPhone() {
	user.loadFromSignInForm(user.AuthType, user.IsoPhone, user.Email)
	user.removeByUserUuid()
}

func (user *User) removeByEmail() {
	user.loadFromSignInForm(user.AuthType, user.IsoPhone, user.Email)
	user.removeByUserUuid()
}
