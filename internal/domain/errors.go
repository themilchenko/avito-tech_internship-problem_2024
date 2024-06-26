package domain

import "errors"

var (
	ErrConflict           = errors.New("conflict")
	ErrBadRequest         = errors.New("bad request")
	ErrUserAlreadyExist   = errors.New("username is already exist")
	ErrBannerAlreadyExist = errors.New(
		"banner with this combination of tagID and featureID already exists",
	)
	ErrBannerNotActive        = errors.New("this banner is not active right now")
	ErrInvalidLoginOrPassword = errors.New("invalid login or password")
	ErrForbidden              = errors.New("you are not supposed to be here")

	ErrCreate = errors.New("failed to create item")
	ErrUpdate = errors.New("failed to update item")
	ErrDelete = errors.New("failed to delete item")

	ErrResponse  = errors.New("failed to response")
	ErrNotFound  = errors.New("failed to find item")
	ErrNoContent = errors.New("no content was found")

	ErrAuth       = errors.New("failed to authenticate")
	ErrNoSession  = errors.New("no existing session")
	ErrBadSession = errors.New("bad session")

	ErrInternal      = errors.New("server error")
	ErrJSONMarshal   = errors.New("failed to marshal json")
	ErrJSONUnmarshal = errors.New("failed to unmarshal json")
	ErrCopy          = errors.New("failed to copy item")
)

var (
	ErrUsername          = errors.New("username not exists")
	ErrPasswordsNotEqual = errors.New("passwords not the same")
)

var ErrRedisNotFound = errors.New("redis: nil")
