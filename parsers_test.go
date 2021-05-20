package cascade

import (
	"context"
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

		parseCounterHouseDto(context.TODO(), data, func(err error) {
			t.Error(err)
		})
	}
}

func TestParseCounterHouseDtoWithCancel(t *testing.T) {
	var test = "/testdata/responses/counterHouse.json"

	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.Fail()
	}

	path := filepath.Join(filepath.Dir(file), test)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	var testItemTotalCount = parseCounterHouseDto(context.TODO(), data, nil)
	t.Logf("total count = %d", testItemTotalCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count int
	items := ParseCounterHouseDto(ctx, data)

	for item := range items {
		if item.error != nil {
			t.Error(item.error)
		}

		if count == 2 {
			cancel()
		}

		count++
	}

	t.Logf("count = %d", count)

	if count == testItemTotalCount {
		t.Fail()
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

		parseCounterHouseReadingDto(context.TODO(), data, func(err error) {
			t.Error(err)
		})
	}
}

func TestParseCounterHouseReadingDtoWithCancel(t *testing.T) {
	var test = "/testdata/responses/readings200.json"

	_, file, _, ok := runtime.Caller(0)

	if !ok {
		t.Fail()
	}

	path := filepath.Join(filepath.Dir(file), test)

	data, err := ioutil.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	var testItemTotalCount = parseCounterHouseReadingDto(context.TODO(), data, nil)
	t.Logf("total count = %d", testItemTotalCount)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var count int
	items := ParseCounterHouseReadingDto(ctx, data)

	for item := range items {
		if item.error != nil {
			t.Error(item.error)
		}

		if count == 2 {
			cancel()
		}

		count++
	}

	t.Logf("count = %d", count)

	if count == testItemTotalCount {
		t.Fail()
	}
}

func parseCounterHouseDto(ctx context.Context, b []byte, errorCallback func(err error)) int {
	var count int

	items := ParseCounterHouseDto(ctx, b)

	for item := range items {
		if item.error != nil {
			if errorCallback != nil {
				errorCallback(item.error)
			}
		}

		count++
	}

	return count
}

func parseCounterHouseReadingDto(ctx context.Context, b []byte, errorCallback func(err error)) int {
	var count int

	items := ParseCounterHouseReadingDto(ctx, b)

	for item := range items {
		if item.error != nil {
			if errorCallback != nil {
				errorCallback(item.error)
			}
		}

		count++
	}

	return count
}
