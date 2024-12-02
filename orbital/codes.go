package orbital

type Code uint32

const (
	OK              Code = 0
	Canceled        Code = 1
	Unknown         Code = 2
	NotFound        Code = 3
	Unimplemented   Code = 4
	Unauthenticated Code = 5
	Internal        Code = 6
	Unavailable     Code = 7
	InvalidRequest  Code = 8

	// TODO: Add more codes as needed
)
