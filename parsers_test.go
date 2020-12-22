package cascade

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"
)

func TestParseCounterHouseDto(t *testing.T) {
	var cases = [1]string{
		"/testdata/responses/counterHouse.json",
	}

	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.Fail()
	}

	for _, test := range cases {
		path := filepath.Join(filepath.Dir(file), test)

		data, err := ioutil.ReadFile(path)

		if err != nil {
			t.Fatal(err)
		}

		parseCounterHouseDto(data, func(err error) {
			t.Error(err)
		})
	}
}

func TestParseCounterHouseReadingDto(t *testing.T) {
	var cases = [1]string{
		"/testdata/responses/readings200.json",
	}

	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.Fail()
	}

	for _, test := range cases {
		path := filepath.Join(filepath.Dir(file), test)

		data, err := ioutil.ReadFile(path)

		if err != nil {
			t.Fatal(err)
		}

		parseCounterHouseReadingDto(data, func(err error) {
			t.Error(err)
		})
	}
}

func parseCounterHouseDto(b []byte, errorCallback func(err error)) {
	items := ParseCounterHouseDto(b)

	for item := range items {
		if item.error != nil {
			errorCallback(item.error)
		}
	}
}

func parseCounterHouseReadingDto(b []byte, errorCallback func(err error)) {
	items := ParseCounterHouseReadingDto(b)

	for item := range items {
		if item.error != nil {
			errorCallback(item.error)
		}
	}
}
