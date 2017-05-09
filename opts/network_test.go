package opts

import (
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNetworkOptLegacySyntax(t *testing.T) {
	testCases := []struct {
		value    string
		expected []network.Options
	}{
		{
			value: "docknet1",
			expected: []network.Options{
				{
					Target: "docknet1",
				},
			},
		},
	}
	for _, tc := range testCases {
		var network NetworkOpt
		assert.NoError(t, network.Set(tc.value))
		assert.Len(t, network.Value(), len(tc.expected))
		for _, expectedNetConfig := range tc.expected {
			verifyNetworkOpt(t, network.Value(), expectedNetConfig)
		}
	}
}

func TestNetworkOptCompleteSyntax(t *testing.T) {
	testCases := []struct {
		value    string
		expected []network.Options
	}{
		{
			value: "name=docknet1,alias=web,driver-opt=field1=value1",
			expected: []network.Options{
				{
					Target:  "docknet1",
					Aliases: []string{"web"},
					DriverOpt: map[string]string{
						"field1": "value1",
					},
				},
			},
		},
		{
			value: "name=docknet1,alias=web1,alias=web2,driver-opt=field1=value1,driver-opt=field2=value2",
			expected: []network.Options{
				{
					Target:  "docknet1",
					Aliases: []string{"web1", "web2"},
					DriverOpt: map[string]string{
						"field1": "value1",
						"field2": "value2",
					},
				},
			},
		},
		{
			value: "name=docknet1",
			expected: []network.Options{
				{
					Target:  "docknet1",
					Aliases: []string{},
				},
			},
		},
	}
	for _, tc := range testCases {
		var network NetworkOpt
		assert.NoError(t, network.Set(tc.value))
		assert.Len(t, network.Value(), len(tc.expected))
		for _, expectedNetConfig := range tc.expected {
			verifyNetworkOpt(t, network.Value(), expectedNetConfig)
		}
	}
}

func TestNetworkOptInvalidSyntax(t *testing.T) {
	testCases := []struct {
		value         string
		expectedError string
	}{
		{
			value:         "invalidField=docknet1",
			expectedError: "invalid field",
		},
		{
			value:         "network=docknet1,invalid=web",
			expectedError: "invalid field",
		},
		{
			value:         "driver-opt=field1=value1,driver-opt=field2=value2",
			expectedError: "network name/id is not specified",
		},
	}
	for _, tc := range testCases {
		var network NetworkOpt
		testutil.ErrorContains(t, network.Set(tc.value), tc.expectedError)
	}
}

func verifyNetworkOpt(t *testing.T, netConfigs []network.Options, expected network.Options) {
	var contains = false
	for _, netConfig := range netConfigs {
		if netConfig.Target == expected.Target {
			if strings.Join(netConfig.Aliases, ",") == strings.Join(expected.Aliases, ",") {
				contains = true
				break
			}
		}
	}
	if !contains {
		t.Errorf("expected %v to contain %v, did not", netConfigs, expected)
	}
}
