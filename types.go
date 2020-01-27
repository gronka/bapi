package main

import (
	"github.com/gocql/gocql"
)

type ApiResponse struct {
	Info      string                 `json:"i"`
	Body      map[string]interface{} `json:"b"`
	FlashMsgs []Msg                  `json:"flashMsgs,omitempty"`
	Command   Command                `json:"command,omitempty"`
}

type Msg struct {
	Msg      string `json:"msg"`
	Severity string `json:"severity"`
}

type Command struct {
	Kind        string `json:"kind,omitempty"`
	Instruction string `json:"instruction,omitempty"`
}

func (ares *ApiResponse) appendFlash(msg, severity string) {
	ares.FlashMsgs = append(ares.FlashMsgs, Msg{msg, severity})
}

type RequestUsingUserUuid struct {
	UserUuid gocql.UUID `json:"userUuid,omitempty", cql:"user_uuid"`
}

var EmptyUuidString = "00000000-0000-0000-0000-000000000000"
var EmptyBytes = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var FullBytes = [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}
var EmptyUuidBytes = [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var AdminUuidString = "99999999-9999-9999-9999-999999999999"
var AdminUuidBytes = [16]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255}

var OneDayMilli = 86400000
var OneDaySecs = 86400

var P_PUBLIC = 1
var P_USER = 2
var P_PUBLIC_MEM_SESSION = 3

var CONVO_ADMIN = 1
var CONVO_BASIC = 2

// after requesting user must select payment
var APPT_STATUS_10_SELECTING_PAYMENT = 10

// after payment selected, we ask the trainer to accept the request
var APPT_STATUS_20_WAIT_FOR_TRAINER_CONFIRMATION = 20

// trainer can reject the request - we notify the player
var APPT_STATUS_35_REJECTED_BY_TRAINER = 35

// trainer accepts - we notify the player
var APPT_STATUS_30_ACCEPTED_BY_TRAINER = 30

// player can cancel in the future
var APPT_STATUS_40_CANCELED_BY_PLAYER = 40

// trainer can cancel in the future
var APPT_STATUS_45_CANCELED_BY_TRAINER = 45

var NO_JWT = "NO_JWT"
var JWT_INVALID = "JWT_INVALID"

var ACCEPTED string = "ACCEPTED"
var REJECTED string = "REJECTED"
var ERROR string = "ERROR"
var NOTSET = "NOTSET"
var DOES_NOT_EXIST = "DOES_NOT_EXIST"

var INFO_REJECTED string = "REJECTED"
var INFO_ACCEPTED string = "ACCEPTED"
var USER_LOGIN_INVALID string = "USER_LOGIN_INVALID"
var USER_VERIFIED string = "USER_VERIFIED"
var USER_BANNED string = "USER_BANNED"

var TOKEN_REJECTED string = "TOKEN_REJECTED"
var AUTH_IS_EMAIL string = "AUTH_IS_EMAIL"
var AUTH_IS_PHONE string = "AUTH_IS_PHONE"
var NO_AUTH string = "NO_AUTH"

var NONE string = "NONE"
var FINE string = "FINE"
var REDIRECT string = "REDIRECT"
var APP string = "App"
var VERIFY_PHONE string = "VerifyPhone"
var VERIFY_EMAIL string = "VerifyEmail"
var EMPTY_STRING string = ""
var REGENERATE string = "REGENERATE"

var SPAM string = "SPAM"

var REQUEST_PENDING = "REQUEST_PENDING"

var TIER_BRONZE = uint64(10)
var TIER_SILVER = uint64(20)
var TIER_GOLD = uint64(30)
var TIER_DIAMOND = uint64(40)

var TIER_BRONZE_BYTE = []byte{10}
var TIER_SILVER_BYTE = []byte{20}
var TIER_GOLD_BYTE = []byte{30}
var TIER_DIAMOND_BYTE = []byte{40}

var TABLE_NAMES = []string{
	"users",
	"user_locs",
	"shouts",
	"shout_locs",
	"courses",
	"course_locs",
	"trainer_applications",
	"phone_verify",
	"phone_verify_spam",
	"planner_updates",
	"planner_available_times",
	"planner_filled_times",

	"convo_sync",
	"convos_by_time",
	"convo_msgs",
	"convo_uuid_to_convo_id",
	"convo_uuid_to_user_uuid",
}
