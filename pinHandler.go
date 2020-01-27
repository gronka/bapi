package main

import (
	//"encoding/binary"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type PinsModel struct {
	// TODO: make sure the muid is never sent to the user
	//Muid      gocql.UUID `json:"muid"`
	EventUuid gocql.UUID `json:"eventUuid"`
	EndTime   int64      `json:"endTime,omitempty"`
	StartTime int64      `json:"startTime,omitempty"`
	TacName   string     `json:"tacName,omitempty"`
	Title     string     `json:"title,omitempty"`
	Lat       float64    `json:"lat,omitempty"`
	Lng       float64    `json:"lng,omitempty"`
	TzOffset  int64      `json:"tzOffset,omitempty"`
	Rating    uint64     `json:"rating,omitempty"`
}

type PinSearchPayload struct {
	Lat       float64 `json:"lat,omitempty"`
	Lng       float64 `json:"lng,omitempty"`
	LatDelta  float64 `json:"latDelta,omitempty"`
	LngDelta  float64 `json:"lngDelta,omitempty"`
	Timestamp int64   `json:"timestamp,omitempty"`
}

type PinSearchResponse struct {
	Zoom    int         `json:"zoom"`
	Bronze  []PinsModel `json:"bronze"`
	Silver  []PinsModel `json:"silver"`
	Gold    []PinsModel `json:"gold"`
	Diamond []PinsModel `json:"diamond"`
}

type PinSearchRange struct {
	Start gocql.UUID
	End   gocql.UUID
}

type PinSearch struct {
	// Needed for init
	Lat       float64
	Lng       float64
	LatDelta  float64
	LngDelta  float64
	Timestamp int64

	BaseMuid []byte
	Psr      PinSearchResponse
}

// numToSelect == 0 means select all
func MakeSearchRange(baseMuid []byte, tier, numToSelect uint64) PinSearchRange {
	var r PinSearchRange
	var err error
	startBytes := make([]byte, 16)
	copy(startBytes, baseMuid)
	startBytes[9] = byte(tier)

	endBytes := make([]byte, 16)
	copy(endBytes, baseMuid)
	endBytes[9] = byte(tier)

	if numToSelect == 0 {
		copy(endBytes[10:], FullBytes[10:])
		r.Start, err = gocql.UUIDFromBytes(startBytes)
		check(err, "failed to convert muid to UUID")
		r.End, err = gocql.UUIDFromBytes(endBytes)
		check(err, "failed to convert muid to UUID")
	}
	return r
}

func (st *PinSearch) SelectAllFromTier(tier uint64) {
	pinRange := MakeSearchRange(st.BaseMuid, tier, 0)
	query := `SELECT * FROM event_pins WHERE muid >= ? AND muid <= ? ALLOW FILTERING`
	iter := s.Query(query, pinRange.Start, pinRange.End).Iter()

	var buf PinsModel
	for {
		row := map[string]interface{}{
			"event_uuid": &buf.EventUuid,
			"end_time":   &buf.EndTime,
			"start_time": &buf.StartTime,
			"tac_name":   &buf.TacName,
			"title":      &buf.Title,
			"lat":        &buf.Lat,
			"lng":        &buf.Lng,
			"tz_offset":  &buf.TzOffset,
			"rating":     &buf.Rating,
		}
		if !iter.MapScan(row) {
			break
		}
		st.Psr.Bronze = append(st.Psr.Bronze, buf)
	}
	err := iter.Close()
	checkIter(err)
}

func (st *PinSearch) Init() {
	mb := MuidBuilder{
		Lat:       st.Lat,
		Lng:       st.Lng,
		LatDelta:  st.LatDelta,
		LngDelta:  st.LngDelta,
		Timestamp: st.Timestamp,
	}
	mb.MakeBaseMuidForSearch()
	st.BaseMuid = mb.BaseMuid
}

func hPinGroupSearchByRegion(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl PinSearchPayload
	decodeJson(r, &pl)

	pinSearch := PinSearch{
		Lat:       pl.Lat,
		Lng:       pl.Lng,
		LatDelta:  pl.LatDelta,
		LngDelta:  pl.LngDelta,
		Timestamp: pl.Timestamp,
	}
	pinSearch.Init()
	pinSearch.SelectAllFromTier(TIER_BRONZE)
	//pins.SelectAllSilver(prefix)
	//pins.SelectAllGold(prefix)
	//pins.SelectAllDiamond(prefix)

	log.Info(pinSearch)
	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["item"] = pinSearch.Psr
	gibs.encodeResponse(w)
}
