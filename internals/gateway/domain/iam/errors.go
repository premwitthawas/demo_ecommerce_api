package iam

import "errors"

var (
	ErrIAMSubjectIsEmpty = errors.New("IAM subject is empty.")
	ErrIAMUnauthoized    = errors.New("IAM unauthorized.")
	ErrIAMForbidden      = errors.New("IAM forbidden.")
	ErrIAMInternalServer = errors.New("IAM internal server error.")
)
