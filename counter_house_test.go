package cascade

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/vitpelekhaty/httptracer"
)

var flowTestCases = []struct {
	value string
	want  Flow
}{
	{
		value: inFlow,
		want:  FlowDirect,
	},
	{
		value: outFlow,
		want:  FlowReverse,
	},
	{
		value: "outFlow",
		want:  FlowReverse,
	},
	{
		value: "test",
		want:  FlowUnknown,
	},
}

func TestCounterHouseChannelDto_Flow(t *testing.T) {
	var have Flow

	for _, test := range flowTestCases {
		channel := &CounterHouseChannelDto{Type: test.value}
		have = channel.Flow()

		if have != test.want {
			t.Errorf(`Flow("%s") failed: have %v, want %v`, test.value, have, test.want)
		}
	}
}

var resourceTestCases = []struct {
	value    string
	resource Resource
}{
	{
		value:    heat,
		resource: ResourceHeat,
	},
	{
		value:    hotWater,
		resource: ResourceHotWater,
	},
	{
		value:    "hotwater",
		resource: ResourceUnknown,
	},
	{
		value:    "test",
		resource: ResourceUnknown,
	},
}

func TestCounterHouseChannelDto_Resource(t *testing.T) {
	var have Resource

	for _, test := range resourceTestCases {
		channel := CounterHouseChannelDto{ResourceType: test.value}
		have = channel.Resource()

		if have != test.resource {
			t.Errorf(`Resource("%s") failed: have %v, want %v`, test.value, have, test.resource)
		}
	}
}

func TestConnection_CounterHouse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockServerFunc))
	defer ts.Close()

	done := false

	_, exec, _, ok := runtime.Caller(0)

	if !ok {
		t.FailNow()
	}

	tracedata := filepath.Join(filepath.Dir(exec), "/testdata/trace/counterHouse_test")

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

	done = true

	data, err := conn.CounterHouse()

	if err != nil {
		t.Fatal(err)
	}

	if len(data) == 0 {
		t.Error("CounterHouse() failed")
	}

	var devices []CounterHouseDto

	err = json.Unmarshal(data, &devices)

	if err != nil {
		t.Fatal(err)
	}
}
