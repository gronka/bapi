package main

import (
	"regexp"
	"strings"

	"github.com/gocql/gocql"
	//"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: allow account hiding/disabling

type User struct {
	UserUuid gocql.UUID `json:"userUuid,omitempty", cql:"user_uuid"`
	Language string     `json:"language,omitempty"`
	Password string     `json:"password,omitempty"`

	Email           string `json:"email,omitempty"`
	EmailVerified   bool   `json:"emailVerified,omitempty", cql:"email_verified"`
	CallingCode     string `json:"callingCode,omitempty", cql:"calling_code"`
	IsoPhone        string `json:"isoPhone,omitempty", cql:"iso_phone"`
	NakedPhone      string `json:"nakedPhone,omitempty"`
	PhoneCountryIso string `json:"phoneCountryIso,omitempty", cql:"phone_country_iso"`
	PhoneVerified   bool   `json:"phoneVerified,omitempty", cql:"phone_verified"`

	Fullname      string `json:"fullname,omitempty"`
	Status        string `json:"status,omitempty"`
	About         string `json:"about,omitempty"`
	Hometown      string `json:"hometown,omitempty"`
	ProfilePicUrl string `json:"profilePicUrl,omitempty", cql:"profile_pic_url"`

	AccountStatus  string  `json:"accountStatus,omitempty", cql:"account_status"`
	LastOnline     uint64  `json:"lastOnline,omitempty", cql:"last_online"`
	OnlineStatus   string  `json:"onlineStatus,omitempty", cql:"online_status"`
	Rating         float32 `json:"rating,omitempty"`
	IsSeller       bool    `json:"isSeller,omitempty", cql:"is_seller"`
	IsTrainer      bool    `json:"isTrainer,omitempty", cql:"is_trainer"`
	IsVideographer bool    `json:"isVideographer,omitempty", cql:"is_videographer"`

	Geohash string  `json:"geohash,omitempty"`
	Lat     float32 `json:"lat,omitempty"`
	Lng     float32 `json:"lng,omitempty"`

	// Never sent in json
	PhoneOrEmail string `json:"phoneOrEmail,omitempty"`
	AuthType     string `json:"authType,omitempty"`
}

func UserFromUserUuid(userUuid gocql.UUID) (user User) {
	user.loadFromUserUuid(userUuid)
	return
}

func UserFromSignInForm(authType, isoPhone, email string) (user User) {
	user.loadFromSignInForm(authType, isoPhone, email)
	return
}

func (user *User) determineAuthType() {
	if len(user.PhoneOrEmail) == 0 {
		user.AuthType = NO_AUTH
		return
	}

	i := strings.Index(user.PhoneOrEmail, "@")
	if i > -1 {
		user.AuthType = AUTH_IS_EMAIL
	} else {
		user.AuthType = AUTH_IS_PHONE
	}
}

func (user *User) prepForSignInAndSignUp() {
	// TODO: put this into User.Unmarshal?
	// TODO: probably remove non-numeric characters from phone numbers

	user.determineAuthType()
	if user.AuthType == AUTH_IS_PHONE {
		// strip non-numbers from phone
		reg := regexp.MustCompile(`[^0-9]`)
		user.NakedPhone = reg.ReplaceAllString(user.PhoneOrEmail, "")
		user.IsoPhone = user.PhoneCountryIso + user.NakedPhone
		// clear these unverified values
		user.Email = ""
	} else if user.AuthType == AUTH_IS_EMAIL {
		// clear these unverified values
		user.NakedPhone = ""
		user.IsoPhone = ""
		user.CallingCode = ""
		user.PhoneCountryIso = ""
		user.Email = strings.TrimSpace(user.PhoneOrEmail)
	} else {
		user.NakedPhone = strings.TrimSpace(user.NakedPhone)
		user.Email = strings.TrimSpace(user.Email)
	}

	user.CallingCode = strings.TrimSpace(user.CallingCode)
	user.IsoPhone = strings.TrimSpace(user.IsoPhone)
	user.PhoneCountryIso = strings.TrimSpace(user.PhoneCountryIso)
	user.Password = strings.TrimSpace(user.Password)
	user.PhoneOrEmail = strings.TrimSpace(user.PhoneOrEmail)

	//printInterface(user)
	log.Info(user.AuthType)
	log.Info(user.IsoPhone)
	log.Info(user.NakedPhone)
}

func (user *User) updateLocation(geohash string, lat, lng float32) {
	user.removeLocs()
	user.saveLatLngGeo(lat, lng, geohash)
	user.saveLoc(geohash[:3], lat, lng)
	user.saveLoc(geohash[:4], lat, lng)
	user.saveLoc(geohash[:5], lat, lng)

	// TODO: check if user is near a disc golf course. If they are, ask them to check in
}
