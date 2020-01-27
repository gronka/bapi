package main

import (
	//"hash"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//log "github.com/sirupsen/logrus"
)

type GetConvoMsgsPayload struct {
	ConvoPayload
}

type GetRecentConvosPayload struct {
	ConvoPayload
}

type Convo struct {
	UserUuid    gocql.UUID `json:"userUuid,omitempty", cql:"user_uuid"`
	ConvoUuid   gocql.UUID `json:"convoUuid,omitempty"`
	LastMsgTime int64      `json:"lastMsgTime", cql:"last_msg_time"`
	Sync        int        `json:"sync", cql:"sync"`
}

type ConvoItem struct {
	LastMsgTime int64      `json:"lastMsgTime", cql:"last_msg_time"`
	ConvoUuid   gocql.UUID `json:"convoUuid,omitempty"`
}

type ConvoMsg struct {
	MsgId        uint       `json:"msgId"`
	TimeSent     int64      `json:"timeSent"`
	MsgUuid      gocql.UUID `json:"msgUuid,omitempty", cql:"msg_uuid"`
	ApparentUuid gocql.UUID `json:"apparentUuid,omitempty", cql:"apparent_uuid"`
	RealUuid     gocql.UUID `json:"realUuid,omitempty", cql:"real_uuid"`
	Body         string     `json:"body"`
}

// realUuid is the userUuid of the user making the request. It will be the
// same as gibs.UserUuid as is passed explicitly
type ConvoPayload struct {
	ConvoUuid    gocql.UUID `json:"convoUuid,omitempty", cql:"convo_uuid"`
	ApparentUuid gocql.UUID `json:"apparentUuid,omitempty", cql:"apparent_uuid"`
	ThisMsgTime  int64
	// client will send these UUIDs sorted
	ParticipantUuids []gocql.UUID `json:"participantUuids,omitempty"`
	ConvoId          []byte       `json:"convoId"`
}

type CreateConvoPayload struct {
	ConvoPayload
}

type ConvoSendMsgPayload struct {
	ConvoPayload
	MsgUuid        gocql.UUID `json:"msgUuid,omitempty", cql:"msg_uuid"`
	OldLastMsgTime int64      `json:"oldLastMsgTime", cql:"old_last_msg_time"`
	Body           string     `json:"body"`
	//ParticipantHash  hash.Hash32  `json:"participantHash,omitempty"`
}

// TODO: no protection against a user creating a convo including users that do
// not exist

func hCreateConvo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl CreateConvoPayload
	decodeJson(r, &pl)

	pl.generateConvoId()

	if pl.convoIdExists() || !pl.realUserCanActAsApparent(gibs.UserUuid) {
		gibs.ares.Info = REJECTED
	} else {
		pl.ConvoUuid, _ = gocql.RandomUUID()
		pl.ThisMsgTime = timeNowMilli()
		pl.create()
		// TODO: send the first message with a body stating who created the convo
		gibs.ares.Info = ACCEPTED
	}
	// TODO: forward user to conversation screen
	gibs.encodeResponse(w)
}

func hConvoSendMsg(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl ConvoSendMsgPayload
	decodeJson(r, &pl)

	if pl.realUserCanActAsApparent(gibs.UserUuid) && pl.userIsConvoParticipantFromDb() {
		pl.ThisMsgTime = timeNowMilli()
		pl.MsgUuid, _ = gocql.RandomUUID()
		pl.send(gibs.UserUuid)
		gibs.ares.Info = ACCEPTED
	} else {
		gibs.ares.Info = REJECTED
	}

	gibs.encodeResponse(w)
}

func hGetRecentConvos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl GetRecentConvosPayload
	decodeJson(r, &pl)

	if pl.realUserCanActAsApparent(gibs.UserUuid) {
		convos := pl.getRecentConvos()
		gibs.ares.Info = ACCEPTED
		gibs.ares.Body["results"] = convos
	} else {
		gibs.ares.Info = REJECTED
	}

	gibs.encodeResponse(w)
}

func hGetConvoMsgs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl GetConvoMsgsPayload
	decodeJson(r, &pl)
	//printInterface(pl)

	if pl.realUserCanActAsApparent(gibs.UserUuid) &&
		pl.userIsConvoParticipantFromDb() {
		convos := pl.getConvoMsgs()
		gibs.ares.Info = ACCEPTED
		gibs.ares.Body["results"] = convos
	} else {
		gibs.ares.Info = REJECTED
	}

	gibs.encodeResponse(w)
}
