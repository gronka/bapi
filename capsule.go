package main

import (
	"strings"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type Capsule struct {
	ItemUuid gocql.UUID  `json:"itemUuid"`
	Field    string      `json:"field"`
	Value    interface{} `json:"value"`
}

var CqlSimpleStatements = map[string]string{
	"about":    `UPDATE users SET about=? WHERE user_uuid=?`,
	"hometown": `UPDATE users SET hometown=? WHERE user_uuid=?`,
	"language": `UPDATE users SET language=? WHERE user_uuid=?`,
	"status":   `UPDATE users SET status=? WHERE user_uuid=?`,
}

func (gibs *Gibs) SanitizeCapsule(capsule *Capsule) bool {
	switch capsule.Field {
	case "about":
		ok := gibs.SanitizeString(capsule, 0, 800)
		return ok

	case "fullname":
		ok := gibs.SanitizeString(capsule, 3, 144)
		return ok

	case "hometown":
		ok := gibs.SanitizeString(capsule, 3, 120)
		return ok

	case "language":
		ok := gibs.SanitizeString(capsule, 2, 20)
		return ok

	case "password":
		ok := gibs.SanitizeString(capsule, 0, 800)
		return ok

	case "status":
		ok := gibs.SanitizeString(capsule, 0, 144)
		return ok

		//fallthrough
	default:
		log.Info("Could not update field with name: " + capsule.Field)
		return false
	}
}

func (gibs *Gibs) SanitizeString(capsule *Capsule, minLength, maxLength int) bool {
	if str, ok := capsule.Value.(string); ok {
		if len(str) <= minLength {
			gibs.ares.Info = REJECTED
			gibs.ares.appendFlash("Submitted value too short", "error")
			return false
		}
		if len(str) > maxLength {
			gibs.ares.Info = REJECTED
			gibs.ares.appendFlash("Submitted value too long", "error")
			return false
		}

		return true
	} else {
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash("Wrong value type submitted", "error")
		return false
	}
}

func (gibs *Gibs) CqlCommitCapsule(capsule *Capsule) {
	switch capsule.Field {
	case "about":
		fallthrough
	case "hometown":
		fallthrough
	case "language":
		fallthrough
	case "status":
		gibs.CqlSimpleUpdate(capsule)

	case "fullname":
		gibs.CqlUpdateFullname(capsule)

	//case "password":
	// probably run a unique function

	default:
		log.Info("Could not update field with name: " + capsule.Field)

	}

}

func (gibs *Gibs) CqlUpdateFullname(capsule *Capsule) {
	fullname, _ := capsule.Value.(string)
	fullnameLower := strings.ToLower(fullname)
	names := strings.Fields(fullnameLower)
	firstName := names[0]
	middleName := ""
	lastName := ""

	// TODO: make names lowercase

	if len(names) == 2 {
		lastName = names[1]
	}

	if len(names) >= 3 {
		middleName = strings.Join(names[1:(len(names)-1)], " ")
		lastName = names[len(names)-1]
	}

	stmt := `UPDATE users SET fullname=?, first_name=?, middle_name=?, last_name=? WHERE user_uuid=?`
	if err := s.Query(stmt, fullname, firstName, middleName, lastName, capsule.ItemUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to update "+capsule.Field))
	}

	gibs.ares.Info = ACCEPTED
	gibs.ares.appendFlash("Update successful", "info")
}

func (gibs *Gibs) CqlSimpleUpdate(capsule *Capsule) {
	stmt := CqlSimpleStatements[capsule.Field]
	err := s.Query(stmt, capsule.Value, capsule.ItemUuid).Exec()
	check(err, "failed to update "+capsule.Field)

	gibs.ares.Info = ACCEPTED
	gibs.ares.appendFlash("Update successful", "info")
}
