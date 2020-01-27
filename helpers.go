package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"reflect"
	"time"

	//"github.com/gocql/gocql"
	//"github.com/julienschmidt/httprouter"
	maps "github.com/gronka/google-maps-services-go"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func timeNowMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

func check(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg))
	}
}

func checkIter(err error) {
	if err != nil {
		panic(errors.Wrap(err, "failed to close iter"))
	}
}

func decodeJson(r *http.Request, target interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&target)
	if err != nil {
		// TODO: get type of target?
		panic(errors.Wrap(err, "failed to decode json"))
	}
}

func printInterface(obj interface{}) {
	ref := reflect.ValueOf(obj)

	for i := 0; i < ref.Type().NumField(); i++ {
		key := ref.Type().Field(i).Name
		typ := ref.Field(i).Type()
		val := ref.Field(i).Interface()

		s := fmt.Sprintf("%s %s = %v;", key, typ, val)
		log.Info().Msg(s)
	}

	//s := reflect.ValueOf(&obj).Elem()
	//typeOfObj := s.Type()

	//for i := 0; i < s.NumField(); i++ {
	//f := s.Field(i)
	//log.Info(`%s %s = %v\n`, i,
	//typeOfObj.Field(i).Name, f.Type(), f.Interface())
	//}
}

func dumpRequest(r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		panic(errors.Wrap(err, "failed to dump request"))
	}
	logInfoMsg("=======http request dump=======")
	logInfoMsg(requestDump)
}

func mapsacParseLocation(location string, r *maps.TimezoneRequest) {
	if location != "" {
		l, err := maps.ParseLatLng(location)
		check(err, "failed to parselocation")
		r.Location = &l
	} else {
		panic("location is required")
	}
}

func printBytes(bytes []byte) {
	for _, n := range bytes {
		fmt.Printf("%8b", n) // prints 1111111111111101
	}
	fmt.Printf("\n")
}

func logInfoMsg(i interface{}) {
	log.Info().Msg(fmt.Sprintf("%s", i))
}

func logFatalMsg(i interface{}) {
	log.Fatal().Msg(fmt.Sprintf("%s", i))
}
