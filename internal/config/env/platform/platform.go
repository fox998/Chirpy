package platform

import (
	"log"
	"strings"
)

type Type int

const (
	Dev Type = iota
	Other
)

func GetPlatformType(typeStr string) Type {

	switch strings.ToLower(typeStr) {
	case "dev":
		return Dev
	default:
		log.Printf("%s not recognised as platform type\n", typeStr)
		return Other
	}
}
