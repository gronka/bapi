package main

import (
//log "github.com/sirupsen/logrus"
)

func SanitizeString(value interface{}, minLength, maxLength int) (string, string) {
	var str, failMsg string
	ok := false
	if str, ok = value.(string); ok {
		if len(str) <= minLength {
			failMsg = "Submitted value too short"
		}
		if len(str) > maxLength {
			failMsg = "Submitted value too long"
		}
	} else {
		failMsg = "Wrong value type submitted"
	}
	return str, failMsg
}

func SanitizeTime(value interface{}) (int64, string) {
	var failMsg string
	var integer64 int64
	ok := false
	if integer64, ok = value.(int64); !ok {
		// TODO: more validation here?
		failMsg = "Wrong value type submitted"
	}
	return integer64, failMsg
}
