package primitive

import (
	"strconv"

	"github.com/google/uuid"
)

type ID uint64

// String returns the string representation of the ID.
func (u ID) String() string {
	return strconv.FormatUint(uint64(u), 10)
}

// Equal checks if two IDs are equal.
func (u ID) Equal(o ID) bool {
	return u == o
}

// ParseID converts a string to an ID. It returns an error if the string is not a valid unsigned integer.
func ParseID(s string) (ID, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return ID(id), nil
}

type UUID = string

func NewUUID() UUID {
	return uuid.NewString()
}

func ParseUUID(s string) (UUID, error) {
	if _, err := uuid.Parse(s); err != nil {
		return "", err
	}
	return s, nil
}
