package errify

type IError interface {
	Error() string
	Message() string
	Location() string
	Details() string
	Code() int
	Read() string

	SetDetails(details string) IError
	JoinLoc(local string) IError
}
