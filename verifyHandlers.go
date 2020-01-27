package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type VerifyPhoneMap struct {
	UserUuid       gocql.UUID `json:"userUuid", cql:"user_uuid"`
	InstallationId string     `json:"installationId", cql:"installation_id"`
	Code           string     `json:"code"`
	ExpireTime     uint64     `json:"expireTime", cql:"expire_time"`
	Attempts       int        `json:"attempts"`
}

type PhoneDetails struct {
	InstallationId string `json:"installationId", cql:"installation_id"`
}

type Guess struct {
	Guess          string `json:"guess"`
	InstallationId string `json:"installationId", cql:"installation_id"`
}

func hVerifyPhoneRegenCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var phoneDetails PhoneDetails
	decodeJson(r, &phoneDetails)

	generateAndCommitPhoneCode(gibs.UserUuid, phoneDetails.InstallationId)
	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hVerifyPhoneCreateCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: make private since userUuid is required
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var phoneDetails PhoneDetails
	decodeJson(r, &phoneDetails)

	var expireTime int64
	stmt := `SELECT expire_time FROM phone_verify WHERE user_uuid=?`
	iter := s.Query(stmt, gibs.UserUuid).Iter()
	iter.Scan(&expireTime)

	var code string
	if expireTime == 0 || time.Now().UTC().Unix() > expireTime {
		code = generateAndCommitPhoneCode(gibs.UserUuid, phoneDetails.InstallationId)
	}

	gibs.ares.Info = ACCEPTED
	// TODO: minor security issue in replying with the code, but it's
	// very useful for testing for now
	gibs.ares.Body["code"] = code
	gibs.encodeResponse(w)
}

func generateAndCommitPhoneCode(userUuid gocql.UUID, installationId string) string {
	code := fmt.Sprintf("%04v", rand.Intn(9999))
	log.Info("Phone verification code: " + code)
	expireTime := time.Now().UTC().Add(30 * time.Minute).Unix()

	stmt := `INSERT INTO phone_verify (user_uuid, installation_id, code, expire_time, attempts) VALUES (?, ?, ?, ?, ?)`
	if err := s.Query(stmt,
		userUuid,
		installationId,
		code,
		expireTime,
		0).Exec(); err != nil {
		panic(errors.Wrap(err, "error writing phone verify code"))
	}
	return code
}

func scrutinizeInstallationIdSpam(installationId string) string {
	var attempts int
	// TODO: handle overloaded requests from one installation id request
	// TODO: store lastupdated time?
	stmt := `SELECT attempts FROM phone_verify_spam WHERE installation_id=?`
	iter := s.Query(stmt, installationId).Iter()
	iter.Scan(attempts)

	stmt = `UPDATE phone_verify_spam SET attempts=? WHERE installation_id=?`
	if err := s.Query(stmt, attempts+1, installationId).Exec(); err != nil {
		panic(errors.Wrap(err, "failed to update"))
	}

	if attempts > 100 {
		return SPAM
	}

	return NONE
}

func hVerifyPhoneCheckCode(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: make private since userUuid is required
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var guessPayload Guess
	decodeJson(r, &guessPayload)
	guessI, err := strconv.Atoi(guessPayload.Guess)
	if err != nil {
		panic(errors.Wrap(err, "failed Atoi"))
	}

	scrut := scrutinizeInstallationIdSpam(guessPayload.InstallationId)
	if scrut != NONE {
		gibs.ResponseCode = 401
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash(
			"Blocked by server",
			"error")
		gibs.encodeResponse(w)
		return
	}

	var attempts int
	var code string
	stmt := `SELECT code, attempts FROM phone_verify WHERE user_uuid=?`
	iter := s.Query(stmt, gibs.UserUuid).Iter()
	iter.Scan(&code, attempts)

	if attempts > 9 {
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash(
			"Failed phone code test too many times",
			"error")
		gibs.encodeResponse(w)
		return
	}

	codeI, err := strconv.Atoi(code)
	if err != nil {
		panic(errors.Wrap(err, "failed Atoi"))
	}

	log.Info("comparing codes")
	log.Info(guessI)
	log.Info(codeI)
	if guessI != codeI {
		attempts++
		stmt := `UPDATE phone_verify SET attempts=? WHERE user_uuid=?`
		if err := s.Query(stmt, attempts, gibs.UserUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "failed to update"))
		}

		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash(
			"Incorrect verification code",
			"error")
	} else {
		// SUCCESS
		stmt := `DELETE FROM phone_verify WHERE user_uuid=?`
		if err := s.Query(stmt, gibs.UserUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "failed to delete row"))
		}

		stmt = `UPDATE users SET phone_verified=? WHERE user_uuid=?`
		if err := s.Query(stmt, true, gibs.UserUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "failed to set phone verified"))
		}

		gibs.ares.Info = ACCEPTED
		gibs.ares.Command.Kind = REDIRECT
		gibs.ares.Command.Instruction = APP
	}

	gibs.encodeResponse(w)
}
