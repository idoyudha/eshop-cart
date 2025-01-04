package utils

import "github.com/google/uuid"

func IDInSliceUUID(a uuid.UUID, uuids uuid.UUIDs) bool {
	for _, b := range uuids {
		if b == a {
			return true
		}
	}
	return false
}
