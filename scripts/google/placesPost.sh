#!/bin/bash

curl -X POST http://127.0.0.1:9090/v1/places/autocomplete \
	-H "Content-Type: application/json" \
	-d '{"input":"3721","latLngString":"38.9072,-77.0369"}'
