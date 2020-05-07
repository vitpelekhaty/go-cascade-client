package cascade

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
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

var beginAt = time.Date(2020, 4, 2, 1, 0, 0, 0, time.UTC)

func TestConnection_Readings_Hours_422(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockServerFunc))
	defer ts.Close()

	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/readingsHours_422_test")

	f, err := os.Create(tracedata)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		f.WriteString("]")
		f.Close()
	}()

	f.WriteString("[")

	client := &http.Client{Timeout: time.Second * 10}
	client = httptracer.Trace(client, httptracer.WithBodies(true), httptracer.WithWriter(f),
		httptracer.WithCallback(func(entry *httptracer.Entry) {
			if !done {
				if entry != nil {
					f.WriteString(",")
				}
			}
		}))

	conn := NewConnection(ts.URL, client)

	defer func() {
		done = true

		if err := conn.Logout(); err != nil {
			t.Error(err)
		}
	}()

	login, err := URLJoin(ts.URL, Login)

	if err != nil {
		t.Fatal(err)
	}

	err = conn.Login(login, auth)

	if err != nil {
		t.Fatal(err)
	}

	ch, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(ch) == 0 {
		t.Error("CounterHouse() failed!")
	}

	var devices []CounterHouseDto

	err = json.Unmarshal(ch, &devices)

	if err != nil {
		t.Fatal(err)
	}

	done = true

	device := devices[0]

	_, err = conn.Readings(device.ID, HourArchive, beginAt, beginAt.Add(time.Hour*24*8))

	if err != nil {
		errorMessage := err.Error()

		if !strings.Contains(errorMessage, "Domain exception has occurred") {
			t.Errorf("Readings 422  failed: have error %s", errorMessage)
		}
	}
}

func TestConnection_Readings_Hours_200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockServerFunc))
	defer ts.Close()

	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/readingsHours_200_test")

	f, err := os.Create(tracedata)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		f.WriteString("]")
		f.Close()
	}()

	f.WriteString("[")

	client := &http.Client{Timeout: time.Second * 10}
	client = httptracer.Trace(client, httptracer.WithBodies(true), httptracer.WithWriter(f),
		httptracer.WithCallback(func(entry *httptracer.Entry) {
			if !done {
				if entry != nil {
					f.WriteString(",")
				}
			}
		}))

	conn := NewConnection(ts.URL, client)

	defer func() {
		done = true

		if err := conn.Logout(); err != nil {
			t.Error(err)
		}
	}()

	login, err := URLJoin(ts.URL, Login)

	if err != nil {
		t.Fatal(err)
	}

	err = conn.Login(login, auth)

	if err != nil {
		t.Fatal(err)
	}

	ch, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(ch) == 0 {
		t.Error("CounterHouse() failed!")
	}

	var devices []CounterHouseDto

	err = json.Unmarshal(ch, &devices)

	if err != nil {
		t.Fatal(err)
	}

	done = true

	device := devices[0]

	data, err := conn.Readings(device.ID, HourArchive, beginAt, beginAt.Add(time.Hour*24))

	if err != nil {
		t.Fatal(err)
	}

	var archive []CounterHouseReadingDto

	err = json.Unmarshal(data, &archive)

	if err != nil {
		t.Fatal(err)
	}

	if len(archive) == 0 {
		t.Error("Readings(HourArchive) failed: empty archive")
	}
}

var dataArchiveStringCases = []struct {
	archive DataArchive
	want    string
}{
	{
		archive: HourArchive,
		want:    "HourArchive",
	},
	{
		archive: DailyArchive,
		want:    "DailyArchive",
	},
	{
		archive: UnknownArchive,
		want:    "UnknownArchive",
	},
}

func TestDataArchive_String(t *testing.T) {
	var have string

	for _, test := range dataArchiveStringCases {
		have = test.archive.String()

		if have != test.want {
			t.Errorf("DataArchive.String() failed: want %s, have %s", test.want, have)
		}
	}
}
