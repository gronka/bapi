package main

import (
	geohasher "github.com/mmcloughlin/geohash"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func latLngRadiusToGeohashList(lat, lng float32, radius int) []string {
	var NORTH_EAST, SOUTH_EAST, SOUTH_WEST, NORTH_WEST bool

	eye := geohasher.Encode(float64(lat), float64(lng))
	precision := precisionFromRadius(radius)
	subRegion := eye[precision]

	// if precision is 4, we always want to expand Noth or South
	switch subRegion {
	case '0':
		fallthrough
	case '1':
		fallthrough
	case '2':
		fallthrough
	case '3':
		fallthrough
	case '4':
		fallthrough
	case '5':
		fallthrough
	case '6':
		fallthrough
	case '7':
		SOUTH_WEST = true

	case 's':
		fallthrough
	case 't':
		fallthrough
	case 'u':
		fallthrough
	case 'v':
		fallthrough
	case 'w':
		fallthrough
	case 'x':
		fallthrough
	case 'y':
		fallthrough
	case 'z':
		NORTH_EAST = true

	case '8':
		fallthrough
	case '9':
		fallthrough
	case 'b':
		fallthrough
	case 'c':
		fallthrough
	case 'd':
		fallthrough
	case 'e':
		fallthrough
	case 'f':
		fallthrough
	case 'g':
		switch precision {
		case 3:
			fallthrough
		case 5:
			fallthrough
		case 7:
			SOUTH_EAST = true
		case 4:
			fallthrough
		case 6:
			NORTH_WEST = true
		default:
			panic(errors.New("Precision not recognized in geosearch"))
		}

	case 'h':
		fallthrough
	case 'j':
		fallthrough
	case 'k':
		fallthrough
	case 'm':
		fallthrough
	case 'n':
		fallthrough
	case 'p':
		fallthrough
	case 'q':
		fallthrough
	case 'r':
		switch precision {
		case 3:
			fallthrough
		case 5:
			fallthrough
		case 7:
			SOUTH_EAST = true
		case 4:
			fallthrough
		case 6:
			NORTH_WEST = true
		default:
			panic(errors.New("Precision not recognized in geosearch"))
		}

	default:
		panic(errors.New("Region not recognized in geosearch"))
	}

	neighbors := geohasher.Neighbors(eye[:precision])
	var n1, n2, n3 string
	if NORTH_EAST {
		log.Info("north east")
		n1 = neighbors[geohasher.North]
		n2 = neighbors[geohasher.NorthEast]
		n3 = neighbors[geohasher.East]

	} else if SOUTH_EAST {
		log.Info("south east")
		n1 = neighbors[geohasher.East]
		n2 = neighbors[geohasher.SouthEast]
		n3 = neighbors[geohasher.South]

	} else if SOUTH_WEST {
		log.Info("south west")
		n1 = neighbors[geohasher.South]
		n2 = neighbors[geohasher.SouthWest]
		n3 = neighbors[geohasher.West]

	} else if NORTH_WEST {
		log.Info("north west")
		n1 = neighbors[geohasher.West]
		n2 = neighbors[geohasher.NorthWest]
		n3 = neighbors[geohasher.North]
	}

	hashes := []string{
		eye[:precision],
		n1[:precision],
		n2[:precision],
		n3[:precision],
	}
	log.Info(hashes)

	return hashes
}

func precisionFromRadius(radius int) (precision uint) {
	// Note: we use a precision of 6 for geohash
	switch radius {
	case 200:
		precision = 3
	case 40:
		precision = 4
	case 10:
		precision = 5
	default:
		panic(errors.New("unsupported search radius"))
	}
	return
}
