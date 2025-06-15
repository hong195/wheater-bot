package request

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParsingCoordsFromStrings(t *testing.T) {
	testCase := []struct {
		str      string
		expected *Coords
		hasErr   bool
	}{
		{
			str: "69.240562,41.311081",
			expected: &Coords{
				Lat: 69.240562,
				Lon: 41.311081,
			},
			hasErr: false,
		},
		{
			str: "/get_weather_details 54.715424,20.509207",
			expected: &Coords{
				Lat: 54.715424,
				Lon: 20.509207,
			},
			hasErr: false,
		},
		{
			str: "//31 54.715424,20.509207",
			expected: &Coords{
				Lat: 54.715424,
				Lon: 20.509207,
			},
			hasErr: false,
		},
		{
			str: "testtest 54.715424,20.509207",
			expected: &Coords{
				Lat: 54.715424,
				Lon: 20.509207,
			},
			hasErr: false,
		},
		{
			str:      "adf 54.715424,",
			expected: &Coords{},
			hasErr:   true,
		},
		{
			str:      "adf 54.715424",
			expected: &Coords{},
			hasErr:   true,
		},
	}

	for _, tc := range testCase {
		coords, err := ParseCoordinates(tc.str)

		if tc.hasErr {
			assert.NotNil(t, err)
			continue
		}

		assert.Nil(t, err)
		assert.Equal(t, tc.expected, coords)
	}
}
