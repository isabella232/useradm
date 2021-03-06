// Copyright 2020 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
package model

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateNew(t *testing.T) {
	testCases := map[string]struct {
		inUser User

		outErr string
	}{
		"email ok, pass ok": {
			inUser: User{
				Email:    "foo@bar.com",
				Password: "correcthorsebatterystaple",
			},
			outErr: "",
		},
		"email invalid, pass ok": {
			inUser: User{
				Email:    "foobar",
				Password: "correcthorsebatterystaple",
			},
			outErr: "email: foobar does not validate as email;",
		},
		"email ok (+), pass ok": {
			inUser: User{
				Email:    "foobar+org@org.com",
				Password: "correcthorsebatterystaple",
			},
			outErr: "",
		},
		"email invalid(non-ascii), pass ok": {
			inUser: User{
				Email:    "ąę@org.com",
				Password: "correcthorsebatterystaple",
			},
			outErr: "email: ąę@org.com does not validate as ascii;",
		},
		"email ok, pass invalid (empty)": {
			inUser: User{
				Email:    "foo@bar.com",
				Password: "",
			},
			outErr: "password can't be empty",
		},
		"email ok, pass invalid (too short)": {
			inUser: User{
				Email:    "foo@bar.com",
				Password: "asdf",
			},
			outErr: "password too short",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			err := tc.inUser.ValidateNew()

			if tc.outErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.outErr)
			}
		})
	}
}

func TestUserFilterParseForm(t *testing.T) {
	testCases := []struct {
		Name string

		Form url.Values

		Result UserFilter
		Error  error
	}{{
		Name: "ok",

		Form: url.Values{
			"id": []string{"1", "2", "3"},
			"email": []string{
				"user1@acme.io",
				"user2@acme.io",
				"user3@acme.io",
			},
			"created_before": []string{"1234567890"},
			"created_after":  []string{"123456789"},
			"updated_before": []string{"9876543210"},
			"updated_after":  []string{"2345678901"},
		},
		Result: UserFilter{
			ID: []string{"1", "2", "3"},
			Email: []string{
				"user1@acme.io",
				"user2@acme.io",
				"user3@acme.io",
			},
			CreatedBefore: func() *time.Time {
				ret := time.Unix(1234567890, 0)
				return &ret
			}(),
			CreatedAfter: func() *time.Time {
				ret := time.Unix(123456789, 0)
				return &ret
			}(),
			UpdatedBefore: func() *time.Time {
				ret := time.Unix(9876543210, 0)
				return &ret
			}(),
			UpdatedAfter: func() *time.Time {
				ret := time.Unix(2345678901, 0)
				return &ret
			}(),
		},
	}, {
		Name: "error, created_after not an int",

		Form: url.Values{
			"created_after": []string{"foo"},
		},
		Error: errors.New(`invalid form parameter "created_after": ` +
			`strconv.ParseInt: parsing "foo": invalid syntax`),
	}, {
		Name: "error, created_before not an int",

		Form: url.Values{
			"created_before": []string{"foo"},
		},
		Error: errors.New(`invalid form parameter "created_before": ` +
			`strconv.ParseInt: parsing "foo": invalid syntax`),
	}, {
		Name: "error, updated_after not an int",

		Form: url.Values{
			"updated_after": []string{"foo"},
		},
		Error: errors.New(`invalid form parameter "updated_after": ` +
			`strconv.ParseInt: parsing "foo": invalid syntax`),
	}, {
		Name: "error, updated_before not an int",

		Form: url.Values{
			"updated_before": []string{"foo"},
		},
		Error: errors.New(`invalid form parameter "updated_before": ` +
			`strconv.ParseInt: parsing "foo": invalid syntax`),
	}}
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var fltr UserFilter

			err := fltr.ParseForm(tc.Form)
			if tc.Error != nil {
				assert.EqualError(t, err, tc.Error.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Result, fltr)
			}
		})
	}
}
