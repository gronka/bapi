package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	geohasher "github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func hUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	users := getAllUsers()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["users"] = users
	gibs.encodeResponse(w)
}

func hUserSignIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var anon User
	decodeJson(r, &anon)
	anon.prepForSignInAndSignUp()

	userFromDb := UserFromSignInForm(anon.AuthType, anon.IsoPhone, anon.Email)
	printInterface(anon)
	printInterface(userFromDb)
	log.Info(userFromDb.PhoneVerified)

	if userFromDb.Password == anon.Password && anon.Password != "" {
		log.Info("this user is logging in " + userFromDb.UserUuid.String())
		// TODO: add user verified case and remove this one
		if userFromDb.AccountStatus == FINE {
			log.Info("logged in as user " + userFromDb.Email)
			log.Info("logged with userUuid " + userFromDb.UserUuid.String())
			jwt := createJwtOnSignIn(userFromDb)
			gibs.ares.Info = ACCEPTED
			gibs.ares.Body["jwt"] = jwt

			// Determine if user has any actions to perform:
			gibs.ares.Command.Kind = REDIRECT
			if userFromDb.PhoneVerified || userFromDb.EmailVerified {
				gibs.ares.Command.Instruction = APP
			} else if !userFromDb.PhoneVerified && userFromDb.IsoPhone != "" {
				gibs.ares.Command.Instruction = VERIFY_PHONE
			} else if !userFromDb.EmailVerified && userFromDb.Email != "" {
				gibs.ares.Command.Instruction = VERIFY_EMAIL
			}

		} else {
			gibs.ResponseCode = 401
			gibs.ares.Info = REJECTED
			gibs.ares.appendFlash(
				"Sign in error: "+userFromDb.AccountStatus,
				"error")
		}

	} else {
		gibs.ResponseCode = 401
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash(REJECTED, "error")
	}

	gibs.encodeResponse(w)
}

func hUserRemove(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: authentication
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var user User
	decodeJson(r, &user)
	user.prepForSignInAndSignUp()

	status := user.lookupAccountStatus()
	if status == DOES_NOT_EXIST {
		log.Info("account does not exist for removal: " + status)
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash("no user to delete", "error")
		gibs.encodeResponse(w)
		return
	}

	if user.UserUuid.String() == EmptyUuidString {
		user.loadFromSignInForm(user.AuthType, user.IsoPhone, user.Email)
	}
	log.Info("====== removing locs")
	user.removeLocs()
	log.Info("====== locs removed")

	if user.AuthType == AUTH_IS_PHONE {
		user.removeByIsoPhone()
	} else if user.AuthType == AUTH_IS_EMAIL {
		user.removeByEmail()
	} else if user.UserUuid.String() != EmptyUuidString {
		user.removeByUserUuid()
	} else {
		panic(errors.New("could not determine method of user removal"))
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hUserSignUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var user User
	decodeJson(r, &user)
	user.prepForSignInAndSignUp()

	status := user.lookupAccountStatus()
	if status != DOES_NOT_EXIST {
		log.Info("account status error: " + status)
		gibs.ares.Info = REJECTED
		gibs.ares.appendFlash("User already exists.", "error")
		gibs.encodeResponse(w)
		return
	}

	user.UserUuid, _ = gocql.RandomUUID()
	user.AccountStatus = FINE
	user.PhoneVerified = false
	user.EmailVerified = false

	user.saveOnSignUp()

	gibs.ares.Command.Kind = REDIRECT
	if user.AuthType == AUTH_IS_PHONE {
		gibs.ares.Command.Instruction = VERIFY_PHONE
	} else if user.AuthType == AUTH_IS_EMAIL {
		gibs.ares.Command.Instruction = VERIFY_EMAIL
	}

	jwt := createJwtOnSignIn(user)
	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["jwt"] = jwt
	gibs.ares.Body["userUuid"] = user.UserUuid
	gibs.ResponseCode = http.StatusCreated
	gibs.encodeResponse(w)
}

func hUserGet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// TODO: block list can be activated here
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var payload RequestUsingUserUuid
	decodeJson(r, &payload)
	log.Info(payload)

	user := UserFromUserUuid(payload.UserUuid)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["user"] = user
	gibs.encodeResponse(w)
}

func hUserFieldUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var capsule Capsule
	decodeJson(r, &capsule)

	authCheck := gibs.authUserCanEditUser(capsule.ItemUuid)
	if !authCheck {
		gibs.sendAuthError(w)
		return
	}

	if ok := gibs.SanitizeCapsule(&capsule); ok {
		gibs.CqlCommitCapsule(&capsule)
	}

	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}

func hUserUpdateLocation(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var locInfo UserLocationPayload
	decodeJson(r, &locInfo)
	lat := locInfo.Lat
	lng := locInfo.Lng
	geohash := geohasher.Encode(float64(lat), float64(lng))

	var user User
	user.UserUuid = gibs.UserUuid
	user.updateLocation(geohash, lat, lng)
	gibs.ares.Info = ACCEPTED
	gibs.encodeResponse(w)
}
