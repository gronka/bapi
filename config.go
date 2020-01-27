package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	//log "github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go"
)

type ConfigStruct struct {
	Environment  string
	bindTo       string
	stripePrefix string
	jsonLoc      string
	mapsac       string
}

var Conf ConfigStruct
var JsonConf map[string]interface{}

func (conf *ConfigStruct) Init() {
	conf.Environment = "dev"
	conf.bindTo = "0.0.0.0:9090"
	conf.stripePrefix = "http://127.0.0.1:9090/v1/stripe"
	home := os.Getenv("HOME")
	conf.jsonLoc = home + "/projects/confs/bapi/scrt.json"
}

func (conf *ConfigStruct) InitJsonConf() {
	jsonFile, err := os.Open(conf.jsonLoc)
	if err != nil {
		panic(errors.Wrap(err, "failed to open json file"))
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal([]byte(byteValue), &JsonConf)
	var ok bool
	if stripe.Key, ok = JsonConf["str"].(string); !ok {
		panic(errors.New("failed to parse JsonConf['str']"))
	}
	if conf.mapsac, ok = JsonConf["mapsac"].(string); !ok {
		panic(errors.New("failed to parse JsonConf['mapsac']"))
	}
	//log.Info(stripe.Key)
}
