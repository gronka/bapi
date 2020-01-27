package main

import (
	"context"
	//"fmt"
	"net/http"

	"github.com/gocql/gocql"
	googleUuid "github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	//"googlemaps.github.io/maps"
	maps "github.com/gronka/google-maps-services-go"
)

// TODO: Note: mapsac api does not support origin parameter - a PR was opened
// TODO: ensure translation beteween gocql.UUID and MapsacSessionToken is correct

type MapsacPayload struct {
	Input        string `json:"input,omitempty"`
	LatLngString string `json:"latLngString,omitempty"`
	PlaceId      string `json:"placeId,omitempty"`
}

type UserHasMapsacTokenFields struct {
	UserUuid gocql.UUID
	Timeout  int64
	Token    gocql.UUID
	Consumed bool
}

func hMapsacCreateSessionToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl MapsacPayload
	token := pl.CreateSessionToken(gibs.UserUuid)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["sessionToken"] = token
	gibs.encodeResponse(w)
}

// TODO: timeout if too manky tokens are made
func (pl *MapsacPayload) CreateSessionToken(userUuid gocql.UUID) gocql.UUID {
	token, _ := gocql.RandomUUID()
	// we don't know how long these tokens last, so we'll set the timeout to 2 minuts
	timeout := timeNowMilli() + (120 * 1000)
	stmt := `INSERT INTO user_has_mapsac_token (user_uuid, timeout, session_token, consumed) VALUES (?, ?, ?, ?) USING TTL ?`
	err := s.Query(stmt, userUuid, timeout, token, false, OneDaySecs).Exec()
	check(err, "error while inserting mapsac token")
	return token
}

func (cqlObj *UserHasMapsacTokenFields) MarkSessionTokenConsumed() {
	stmt := `UPDATE user_has_mapsac_token SET consumed=? WHERE user_uuid=? AND timeout=? AND session_token=?`
	err := s.Query(stmt, true, cqlObj.UserUuid, cqlObj.Timeout, cqlObj.Token).Exec()
	check(err, "error while setting token consumed true")
}

func (pl *MapsacPayload) GetSessionInfo(userUuid gocql.UUID) UserHasMapsacTokenFields {
	var result UserHasMapsacTokenFields
	stmt := `SELECT * FROM user_has_mapsac_token WHERE user_uuid=? LIMIT 1`
	iter := s.Query(stmt, userUuid).Iter()

	for {
		row := map[string]interface{}{
			"user_uuid":     &result.UserUuid,
			"timeout":       &result.Timeout,
			"session_token": &result.Token,
			"consumed":      &result.Consumed,
		}
		if !iter.MapScan(row) {
			break
		}
	}
	err := iter.Close()
	check(err, "failed to close iter")

	if (timeNowMilli() > result.Timeout) || (result.Consumed) {
		result.Token = pl.CreateSessionToken(userUuid)
	}

	return result
}

//func (pl *PlacesPayload) CreateQuery() {
//base := "https://maps.googleapis.com/maps/api/place/autocomplete/json?key=" + Conf.places + "&input=" + pl.Input
//if pl.Lat != 0 && pl.Lng != 0 {
//latStr := fmt.Sprintf("%f", pl.Lat)
//lngStr := fmt.Sprintf("%f", pl.Lng)
//base += "&location=" + latStr + "," + lngStr
//base += "&origin=" + latStr + "," + lngStr
//}
//}

func SessionTokenFromGocqlUuid(uuid gocql.UUID) maps.PlaceAutocompleteSessionToken {
	asGoogleUuid, _ := googleUuid.FromBytes(uuid.Bytes())
	token := maps.PlaceAutocompleteSessionToken(asGoogleUuid)
	return token
}

func parseLocation(location string, req *maps.PlaceAutocompleteRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err, "failed to parse latlng")
		req.Location = &l
	}
}

func parseOrigin(location string, req *maps.PlaceAutocompleteRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err, "failed to parse latlng")
		req.Origin = &l
	}
}

func hMapsacPredictions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl MapsacPayload
	decodeJson(r, &pl)
	cqlObj := pl.GetSessionInfo(gibs.UserUuid)

	sessionToken := SessionTokenFromGocqlUuid(cqlObj.Token)
	mapsacReq := &maps.PlaceAutocompleteRequest{
		Input:        pl.Input,
		Radius:       500,
		SessionToken: sessionToken,
	}
	parseLocation(pl.LatLngString, mapsacReq)
	parseOrigin(pl.LatLngString, mapsacReq)

	client, err := maps.NewClient(maps.WithAPIKey(Conf.mapsac))
	check(err, "failed to create mapsac client")
	mapsacRes, err := client.PlaceAutocomplete(context.Background(), mapsacReq)
	log.Info(gibs.UserUuid.String() + " made mapsac request")
	check(err, "failed place autocomplete")

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["predictions"] = mapsacRes.Predictions
	gibs.encodeResponse(w)
}

func hMapsacLookupByPlaceId(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var pl MapsacPayload
	decodeJson(r, &pl)

	cqlObj := pl.GetSessionInfo(gibs.UserUuid)
	sessionToken := SessionTokenFromGocqlUuid(cqlObj.Token)
	log.Info("place id")
	log.Info(pl.PlaceId)
	mapsacReq := &maps.PlaceDetailsRequest{
		PlaceID:      pl.PlaceId,
		SessionToken: sessionToken,
	}

	client, err := maps.NewClient(maps.WithAPIKey(Conf.mapsac))
	check(err, "failed to create mapsac client")
	mapsacRes, err := client.PlaceDetails(context.Background(), mapsacReq)
	check(err, "failed place autocomplete")
	cqlObj.MarkSessionTokenConsumed()

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["place"] = mapsacRes
	gibs.encodeResponse(w)
}
