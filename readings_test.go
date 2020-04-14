package cascade

import (
	"encoding/json"
	"testing"
	"time"
)

var requestTimeStringTestCases = []struct {
	value time.Time
	want  string
}{
	{
		value: time.Date(2020, 4, 14, 11, 21, 15, 0, time.UTC),
		want:  "14.04.2020 11:21:15",
	},
}

func TestRequestTime_String(t *testing.T) {
	var (
		rt   RequestTime
		have string
	)

	for _, test := range requestTimeStringTestCases {
		rt = RequestTime(test.value)
		have = rt.String()

		if test.want != rt.String() {
			t.Errorf("%q.String() failed: have %s, want %s", test.value, have, test.want)
		}
	}
}

func TestRequestTime_UnmarshalJSON(t *testing.T) {
	var rt RequestTime

	for _, test := range requestTimeStringTestCases {
		err := rt.UnmarshalJSON([]byte(test.want))

		if err != nil {
			t.Fatal(err)
		}

		if time.Time(rt) != test.value {
			t.Errorf(`RequestTime.UnmarshalJSON("%s") failed: have %q, want %q`, test.want, time.Time(rt),
				test.value)
		}
	}
}

func TestRequestTime_MarshalJSON(t *testing.T) {
	rt := RequestTime(time.Date(2019, 12, 11, 1, 0, 0, 0, time.UTC))

	data, err := json.Marshal(&rt)

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Errorf("RequestTime.Marshal(%v) failed", time.Time(rt))
	}
}

var dataArchiveUnmarshalJSONTestCases = []struct {
	value string
	want  DataArchive
}{
	{
		value: "1",
		want:  HourArchive,
	},
	{
		value: "2",
		want:  DailyArchive,
	},
	{
		value: "0",
		want:  UnknownArchive,
	},
	{
		value: "10",
		want:  UnknownArchive,
	},
	{
		value: "",
		want:  UnknownArchive,
	},
}

func TestDataArchive_UnmarshalJSON(t *testing.T) {
	var archive DataArchive

	for _, test := range dataArchiveUnmarshalJSONTestCases {
		archive.UnmarshalJSON([]byte(test.value))

		if archive != test.want {
			t.Errorf(`DataArchive.UnmarshalJSON("%s") failed: have %q, want %q`, test.value, archive, test.want)
		}
	}
}

var readingsDataArchiveTestCases = []struct {
	value int
	want  DataArchive
}{
	{
		value: 1,
		want:  HourArchive,
	},
	{
		value: 2,
		want:  DailyArchive,
	},
	{
		value: 10,
		want:  UnknownArchive,
	},
}

func TestCounterHouseReadingDto_DataArchive(t *testing.T) {
	var have DataArchive

	for _, test := range readingsDataArchiveTestCases {
		r := CounterHouseReadingDto{
			Archive: int32(test.value),
		}

		have = r.DataArchive()

		if have != test.want {
			t.Errorf("CounterHouseReadingDto.DataArchive(%d) failed: have %v, want %v", test.value, have,
				test.want)
		}
	}
}
