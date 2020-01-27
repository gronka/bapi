package main

import (
	"net/http"
	//"strings"

	"github.com/gocql/gocql"
	"github.com/julienschmidt/httprouter"
	//log "github.com/sirupsen/logrus"
)

type CourseSearchResult struct {
	CourseUuid    gocql.UUID `json:"courseUuid,omitempty", cql:"course_uuid"`
	Fullname      string     `json:"fullname,omitempty"`
	Popularity    float32    `json:"popularity,omitempty"`
	Randomize     float32    `json:"randomize,omitempty"`
	ProfilePicUrl string     `json:"profilePicUrl,omitempty", cql:"is_videographer"`
	Status        string     `json:"status,omitempty"`
	SearchScore   int        `json:"searchScore"`

	Lat float32 `json:"lat,omitempty"`
	Lng float32 `json:"lng,omitempty"`
}

type SearchCoursesPayload struct {
	Lat     float32              `json:"lat"`
	Lng     float32              `json:"lng"`
	Radius  int                  `json:"radius"`
	Name    string               `json:"name"`
	Results []CourseSearchResult `json:"results"`
}

func hSearchCoursesByNearest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gibs := unpackReqIntoGibs(w, r, P_PUBLIC)
	var scp SearchCoursesPayload
	decodeJson(r, &scp)

	geohashes := latLngRadiusToGeohashList(scp.Lat, scp.Lng, scp.Radius)
	scp.searchByGeohashes(geohashes)

	gibs.ares.Info = ACCEPTED
	gibs.ares.Body["results"] = scp.Results
	gibs.encodeResponse(w)
}
