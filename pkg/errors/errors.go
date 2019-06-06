package errors

import "errors"

// Common errors
var (
	ErrDatabase                        = errors.New("something wrong with database")
	ErrExchange                        = errors.New("something wrong with exchange")
	ErrUnexpected                      = errors.New("unexpected error")
	ErrUserIsNotFound                  = errors.New("user is not found")
	ErrUserHasNotAuthToken             = errors.New("user has not auth token")
	ErrUserHasWrongAuthToken           = errors.New("user has wrong auth token")
	ErrUserIsNotAuthenticatedWithToken = errors.New("need complete authentication with token")
	ErrUserWithEmailNotFound           = errors.New("user with email not found")
	ErrUserEmailIsAlreadyRegistered    = errors.New("user email is already registered")
	ErrUserIsAlreadyAuthenticated      = errors.New("user is already authenticated")
	ErrThereAreNotMeetings             = errors.New("user has not meetings")
	ErrUserHasNotMeeting               = errors.New("user has not meeting")
	ErrUserWithoutRoomlist             = errors.New("user without roomList")
)
