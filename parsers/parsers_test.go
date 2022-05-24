package parsers

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGaugesList(t *testing.T) {
	var cases = [1]string{
		"../testdata/responses/counterHouse.json",
	}

	var err error

	for _, test := range cases {
		path := test

		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(path)

			require.NoError(t, err, path)
		}

		data, err := ioutil.ReadFile(path)

		require.NoError(t, err, path)

		var row int

		items, err := ParseGaugesList(context.TODO(), data)

		require.NoError(t, err, path)

		for item := range items {
			assert.NoError(t, item.E, path, row)

			if !item.Error() {
				_, ok := item.V.(*Gauge)

				assert.Equal(t, true, ok, path, row)
			}

			row++
		}
	}
}

func TestParseReadings(t *testing.T) {
	var cases = [1]string{
		"../testdata/responses/readings200.json",
	}

	var err error

	for _, test := range cases {
		path := test

		if !filepath.IsAbs(path) {
			path, err = filepath.Abs(path)

			require.NoError(t, err, path)
		}

		data, err := ioutil.ReadFile(path)

		require.NoError(t, err, path)

		var row int

		items, err := ParseReadings(context.TODO(), data)

		require.NoError(t, err, path)

		for item := range items {
			assert.NoError(t, item.E, path, row)

			if !item.Error() {
				_, ok := item.V.(*Readings)

				assert.Equal(t, true, ok, path, row)
			}

			row++
		}
	}
}
