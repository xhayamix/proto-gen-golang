// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: client/api/user_gen.proto

package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on Profile with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Profile) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Profile with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ProfileMultiError, or nil if none found.
func (m *Profile) ValidateAll() error {
	return m.validate(true)
}

func (m *Profile) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for AccountId

	// no validation rules for Name

	// no validation rules for IconImageUrl

	// no validation rules for Bio

	if len(errors) > 0 {
		return ProfileMultiError(errors)
	}

	return nil
}

// ProfileMultiError is an error wrapping multiple validation errors returned
// by Profile.ValidateAll() if the designated constraints aren't met.
type ProfileMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ProfileMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ProfileMultiError) AllErrors() []error { return m }

// ProfileValidationError is the validation error returned by Profile.Validate
// if the designated constraints aren't met.
type ProfileValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ProfileValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ProfileValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ProfileValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ProfileValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ProfileValidationError) ErrorName() string { return "ProfileValidationError" }

// Error satisfies the builtin error interface
func (e ProfileValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sProfile.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ProfileValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ProfileValidationError{}

// Validate checks the field values on GetProfileRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetProfileRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetProfileRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetProfileRequestMultiError, or nil if none found.
func (m *GetProfileRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetProfileRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetUserId()) < 1 {
		err := GetProfileRequestValidationError{
			field:  "UserId",
			reason: "value length must be at least 1 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return GetProfileRequestMultiError(errors)
	}

	return nil
}

// GetProfileRequestMultiError is an error wrapping multiple validation errors
// returned by GetProfileRequest.ValidateAll() if the designated constraints
// aren't met.
type GetProfileRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetProfileRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetProfileRequestMultiError) AllErrors() []error { return m }

// GetProfileRequestValidationError is the validation error returned by
// GetProfileRequest.Validate if the designated constraints aren't met.
type GetProfileRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetProfileRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetProfileRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetProfileRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetProfileRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetProfileRequestValidationError) ErrorName() string {
	return "GetProfileRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetProfileRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetProfileRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetProfileRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetProfileRequestValidationError{}

// Validate checks the field values on GetProfileResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetProfileResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetProfileResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetProfileResponseMultiError, or nil if none found.
func (m *GetProfileResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetProfileResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetProfile()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetProfileResponseValidationError{
					field:  "Profile",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetProfileResponseValidationError{
					field:  "Profile",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetProfile()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetProfileResponseValidationError{
				field:  "Profile",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetProfileResponseMultiError(errors)
	}

	return nil
}

// GetProfileResponseMultiError is an error wrapping multiple validation errors
// returned by GetProfileResponse.ValidateAll() if the designated constraints
// aren't met.
type GetProfileResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetProfileResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetProfileResponseMultiError) AllErrors() []error { return m }

// GetProfileResponseValidationError is the validation error returned by
// GetProfileResponse.Validate if the designated constraints aren't met.
type GetProfileResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetProfileResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetProfileResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetProfileResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetProfileResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetProfileResponseValidationError) ErrorName() string {
	return "GetProfileResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetProfileResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetProfileResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetProfileResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetProfileResponseValidationError{}