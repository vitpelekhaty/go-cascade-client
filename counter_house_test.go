package cascade

import (
	"testing"
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
		want:  FlowUnknown,
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
