package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/handlers"
	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var s *gocql.Session

func init() {
	// TODO: load info from a config file
	Conf.Init()
	Conf.InitJsonConf()
}

func main() {
	//log.SetLevel(log.InfoLevel)
	// TODO: refactor to use zerolog
	// TODO: set production flag to disable pretty logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	reinit := flag.Bool("reinit", false, "does nothing rn")

	casshost := "127.0.0.1"
	//casshost := "localhost"

	cluster := gocql.NewCluster(casshost)
	cluster.Keyspace = "bapi"
	cluster.Consistency = gocql.One
	cluster.ProtoVersion = 3
	var err error
	s, err = cluster.CreateSession()
	for err != nil {
		s, err = cluster.CreateSession()
		time.Sleep(1 * time.Second)
		log.Error().Msg("CQL session could not be established. Retrying in 1 second.")
	}
	log.Info().Msg("CQL session established! Starting server!")
	defer s.Close()

	if *reinit == true {
		// TODO: other commands

	} else {
		r := ApiRouter()

		headersOk := handlers.AllowedHeaders([]string{
			"Authorization",
			"Content-Type",
			"X-Requested-With",
			"X-Session-Token",
		})
		originsOk := handlers.AllowedHeaders([]string{"*"})
		methodsOk := handlers.AllowedHeaders([]string{
			"GET",
			"POST",
			"PUT",
			"OPTIONS",
			"DELETE",
		})

		//log.Fatal(http.ListenAndServe(":9090", handlers.CORS(originsOk, headersOk, methodsOk)(r)))

		srv := &http.Server{
			//Handler:      handlers.LoggingHandler(os.Stdout, r),
			Handler: handlers.LoggingHandler(os.Stdout,
				handlers.CORS(originsOk, headersOk, methodsOk)(r)),
			Addr:         Conf.bindTo,
			WriteTimeout: 15 * time.Second,
			ReadTimeout:  15 * time.Second,
		}

		// TODO: check that all errors are wrapped before coming here
		//log.Fatal().Msg(errors.Cause(srv.ListenAndServe()))
		log.Fatal().Msg(fmt.Sprintf("%s", errors.Cause(srv.ListenAndServe())))
	}
}
