package repository

import "errors"

var (
	//ErrEventExists returned when there is an attemp to create new event with same uid
	ErrEventExists = errors.New("event with this uid already exists")
	//ErrEventDoesNotExists ...
	ErrEventDoesNotExists = errors.New("event does not exists")
	//ErrRequiredFieldsNotProvided ...
	ErrRequiredFieldsNotProvided = errors.New("required fields are not provided")
	//ErrIncorrectValue ...
	ErrIncorrectValue = errors.New("incorrect value provided for the field")
)
