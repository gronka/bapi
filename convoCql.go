package main

import (
	"bytes"

	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// generateConvoId creates the ConvoId from the supplied ConvoUuid and the
// UUID of all chat participants. The purpose of the ConvoId is to optimize
// and secure: it removes the need to SELECT all convo ParticipantUuids to send
// a message to; it ensures that the list of ParticipantsUuids supplied by the
// client is not spoofed.
func (pl *CreateConvoPayload) generateConvoId() {
	// assume the given uuids are already sorted
	var b bytes.Buffer
	// TODO: only add convouuid to the convoId when there is more than 2 convo
	// participants. This allows groups of people to have multiple conversations
	// while restricting one on one convos to just two people
	for _, byt := range pl.ConvoUuid.Bytes() {
		b.WriteByte(byt)
	}
	for _, uuid := range pl.ParticipantUuids {
		for _, byt := range uuid.Bytes() {
			b.WriteByte(byt)
		}
	}
	pl.ConvoId = b.Bytes()
}

func (pl *CreateConvoPayload) convoIdExists() bool {
	iter := s.Query(`SELECT convo_uuid FROM convo_uuid_to_convo_id
	WHERE convo_id = ?`, pl.ConvoId).Iter()
	iter.Scan(
		&pl.ConvoUuid,
	)

	if pl.ConvoUuid == EmptyUuidBytes {
		return false
	}
	return true
}

// userIsConvoParticipant checks that the user is a member of the convo. This
// function only works after the ConvoId generated from the participant list
// has been proven valid. Specifically, after generateConvoId and convoExists
// have been executed
func (pl *ConvoPayload) userIsConvoParticipantFromConvoId(gibs Gibs) bool {
	// TODO: check blacklist of convo
	for _, uuid := range pl.ParticipantUuids {
		if gibs.UserUuid == uuid {
			return true
		}
	}
	return false
}

func (pl *ConvoPayload) userIsConvoParticipantFromDb() bool {
	var convoUuidLookup gocql.UUID
	iter := s.Query(`SELECT convo_uuid FROM convo_uuid_to_user_uuid
WHERE convo_uuid = ? AND user_uuid = ?`, pl.ConvoUuid, pl.ApparentUuid).Iter()
	iter.Scan(
		&convoUuidLookup,
	)
	if convoUuidLookup == EmptyUuidBytes {
		return false
	}
	return true
}

// realUserCanActAsApparent checks that the logged in user has permission to
// send the message with the supplied ApparentUuid.
func (pl *ConvoPayload) realUserCanActAsApparent(realUuid gocql.UUID) bool {
	// TODO: if orgs are ever added, we could check here to see if the
	// user is allowed to represent the org

	// TODO: check blacklist

	// TODO: allow admin accounts here

	if realUuid == pl.ApparentUuid {
		return true
	}
	return false
}

func (pl *CreateConvoPayload) create() {
	log.Info(pl.ConvoId)
	log.Info(pl.ConvoUuid)
	stmt := `INSERT INTO convo_uuid_to_convo_id (convo_id, convo_uuid) VALUES (?, ?)`
	if err := s.Query(stmt,
		pl.ConvoId,
		pl.ConvoUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

	for _, participantUuid := range pl.ParticipantUuids {
		// Set last_msg_time for all participants
		stmt = `INSERT INTO convos_by_time (user_uuid, convo_uuid, last_msg_time) VALUES (?, ?, ?)`
		if err := s.Query(stmt,
			participantUuid,
			pl.ConvoUuid,
			pl.ThisMsgTime).Exec(); err != nil {
			panic(errors.Wrap(err, "error while querying"))
		}

		// Add users to convo list
		stmt = `INSERT INTO convo_uuid_to_user_uuid (convo_uuid, user_uuid) VALUES (?, ?)`
		if err := s.Query(stmt, pl.ConvoUuid, participantUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "error while querying"))
		}

	}
}

func getIncrementalMsgId(convoUuid gocql.UUID) uint {
	var msgId uint
	stmt := `SELECT msg_id FROM convo_msg_id_counter WHERE convo_uuid=?`
	iter := s.Query(stmt, convoUuid).Iter()
	iter.Scan(&msgId)

	stmt = `UPDATE convo_msg_id_counter SET msg_id = msg_id + 1 WHERE convo_uuid = ?`
	if err := s.Query(stmt, convoUuid).Exec(); err != nil {
		panic(errors.Wrap(err, "error while incrementing msg_id"))
	}
	return msgId
}

func (pl *ConvoSendMsgPayload) send(realUuid gocql.UUID) {
	msgId := getIncrementalMsgId(pl.ConvoUuid)

	stmt := `INSERT INTO convo_msgs (convo_uuid, msg_id, time_sent, msg_uuid, apparent_uuid, real_uuid, body) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := s.Query(stmt,
		pl.ConvoUuid,
		msgId,
		pl.ThisMsgTime,
		pl.MsgUuid,
		pl.ApparentUuid,
		realUuid,
		pl.Body).Exec(); err != nil {
		panic(errors.Wrap(err, "error while querying"))
	}

	// TODO: look up participantuuids
	for _, participantUuid := range pl.ParticipantUuids {
		// Set last_msg_time for all participants
		stmt = `UPDATE convos_by_time SET last_msg_time = ? WHERE user_uuid = ? AND last_msg_time = ? AND convo_uuid = ?`
		if err := s.Query(stmt,
			pl.ThisMsgTime,
			participantUuid,
			pl.OldLastMsgTime,
			pl.ConvoUuid).Exec(); err != nil {
			panic(errors.Wrap(err, "error while querying"))
		}
	}
}

func (pl *GetRecentConvosPayload) getRecentConvos() []ConvoItem {
	// TODO: add paging
	iter := s.Query(`SELECT last_msg_time, convo_uuid FROM convos_by_time
	WHERE user_uuid = ?`, pl.ApparentUuid).Iter()
	var convos []ConvoItem
	var convo ConvoItem
	for {
		row := map[string]interface{}{
			"last_msg_time": &convo.LastMsgTime,
			"convo_uuid":    &convo.ConvoUuid,
		}
		if !iter.MapScan(row) {
			break
		}
		convos = append(convos, convo)
	}
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
	//for {
	//iter.Scan(
	//&convo.LastMsgTime,
	//&convo.ConvoUuid,
	//)
	//convos = append(convos, convo)
	//}

	return convos
}

func (pl *GetConvoMsgsPayload) getConvoMsgs() []ConvoMsg {
	// TODO: should this be getRecentConvoMsgs?
	// TODO: add paging
	// TODO: using the lastMsgTime on the client device for getting new
	// messages would be more efficient and effective than using convo_sync
	// TODO: authorize user
	// TODO: make sure user is member of conversation and has okay status
	iter := s.Query(`SELECT time_sent, msg_uuid, apparent_uuid, real_uuid, body FROM convo_msgs
	WHERE convo_uuid = ?`, pl.ConvoUuid).Iter()
	var msgs []ConvoMsg
	var msg ConvoMsg
	for {
		row := map[string]interface{}{
			"msg_id":        &msg.MsgId,
			"time_sent":     &msg.TimeSent,
			"msg_uuid":      &msg.MsgUuid,
			"apparent_uuid": &msg.ApparentUuid,
			"real_uuid":     &msg.RealUuid,
			"body":          &msg.Body,
		}
		if !iter.MapScan(row) {
			break
		}
		msgs = append(msgs, msg)
	}
	if err := iter.Close(); err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}

	//for {
	//iter.Scan(
	//&msg.TimeSent,
	//&msg.MsgUuid,
	//&msg.ApparentUuid,
	//&msg.Body,
	//)
	//msgs = append(msgs, msg)
	//}

	return msgs
}

//func (pl *CreateConvoPayload) generateParticipantHash() {
//// Note: this function will need to be kept current in relation to how
//// conversations are stored
//uuidCount := len(pl.ParticipantUuids)
//listToHash := make([]string, uuidCount*2)
//for i, uuid := range pl.ParticipantUuids {
//thisItem := uuid.String() + string(pl.ParticipantRoles[i])
//listToHash = append(listToHash, thisItem)
//}

//sort.Strings(listToHash)
//asString := strings.Join(listToHash, "")

//hash := fnv.New64()
//hash.Write([]byte(asString))
//log.Info(hash.Sum64())
//value := hash.Sum64()

//b := make([]byte, 16)
//binary.LittleEndian.PutUint64(b, value)

//pl.ParticipantHash, _ = gocql.UUIDFromBytes(b)
