package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAzSegment(t *testing.T) {
	cases := []struct {
		Case            string
		ExpectedEnabled bool
		ExpectedString  string
		EnvSubName      string
		EnvSubID        string
		CliExists       bool
		CliSubName      string
		CliSubID        string
		InfoSeparator   string
		DisplayID       bool
		DisplayName     bool
	}{
		{Case: "envvars present",
			ExpectedEnabled: true,
			ExpectedString:  "foo$bar",
			EnvSubName:      "foo",
			EnvSubID:        "bar",
			CliExists:       false,
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     true},
		{Case: "envvar name present",
			ExpectedEnabled: true,
			ExpectedString:  "foo$",
			EnvSubName:      "foo",
			EnvSubID:        "",
			CliExists:       false,
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     true},
		{Case: "envvar id present",
			ExpectedEnabled: true,
			ExpectedString:  "$bar",
			EnvSubName:      "",
			EnvSubID:        "bar",
			CliExists:       false,
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     true},
		{Case: "cli not found",
			ExpectedEnabled: false,
			ExpectedString:  "$",
			EnvSubName:      "",
			EnvSubID:        "",
			CliExists:       false,
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     true},
		{Case: "cli contains data",
			ExpectedEnabled: true,
			ExpectedString:  "foo$bar",
			EnvSubName:      "",
			EnvSubID:        "",
			CliExists:       true,
			CliSubName:      "foo",
			CliSubID:        "bar",
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     true},
		{Case: "print only name",
			ExpectedEnabled: true,
			ExpectedString:  "foo",
			EnvSubName:      "",
			EnvSubID:        "",
			CliExists:       true,
			CliSubName:      "foo",
			CliSubID:        "bar",
			InfoSeparator:   "$",
			DisplayID:       false,
			DisplayName:     true},
		{Case: "print only id",
			ExpectedEnabled: true,
			ExpectedString:  "bar",
			EnvSubName:      "",
			EnvSubID:        "",
			CliExists:       true,
			CliSubName:      "foo",
			CliSubID:        "bar",
			InfoSeparator:   "$",
			DisplayID:       true,
			DisplayName:     false},
		{Case: "print none",
			ExpectedEnabled: false,
			ExpectedString:  "",
			EnvSubName:      "",
			EnvSubID:        "",
			CliExists:       true,
			CliSubName:      "foo",
			CliSubID:        "bar",
			InfoSeparator:   "$",
			DisplayID:       false,
			DisplayName:     false},
	}

	for _, tc := range cases {
		env := new(MockedEnvironment)
		env.On("getenv", "AZ_SUBSCRIPTION_NAME").Return(tc.EnvSubName)
		env.On("getenv", "AZ_SUBSCRIPTION_ID").Return(tc.EnvSubID)
		env.On("hasCommand", "az").Return(tc.CliExists)
		env.On("runCommand", "az", []string{"account", "show", "--query=[name,id]", "-o=tsv"}).Return(fmt.Sprintf("%s\n%s\n", tc.CliSubName, tc.CliSubID), nil)
		props := &properties{
			values: map[Property]interface{}{
				SubscriptionInfoSeparator: tc.InfoSeparator,
				DisplaySubscriptionID:     tc.DisplayID,
				DisplaySubscriptionName:   tc.DisplayName,
			},
		}

		az := &az{
			env:   env,
			props: props,
		}
		assert.Equal(t, tc.ExpectedEnabled, az.enabled(), tc.Case)
		assert.Equal(t, tc.ExpectedString, az.string(), tc.Case)
	}
}
