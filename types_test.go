package unifi_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unpoller/unifi/v5"
)

func TestFlexInt(t *testing.T) {
	t.Parallel()
	a := assert.New(t)
	five, seven := 5, 7

	var r struct {
		Five    unifi.FlexInt `json:"five"`
		Seven   unifi.FlexInt `json:"seven"`
		Auto    unifi.FlexInt `json:"auto"`
		Channel unifi.FlexInt `json:"channel"`
		Nil     unifi.FlexInt `json:"nil"`
	}

	// test unmarshalling the custom type three times with different values.
	a.Nil(json.Unmarshal([]byte(`{"five": "5", "seven": 7, "auto": "auto", "nil": null}`), &r))
	// test number in string.
	a.EqualValues(five, r.Five.Val)
	a.EqualValues("5", r.Five.Txt)
	// test number.
	a.EqualValues(seven, r.Seven.Val)
	a.EqualValues("7", r.Seven.Txt)
	// test string.
	a.EqualValues(0, r.Auto.Val)
	a.EqualValues("auto", r.Auto.Txt)
	// test (error) struct.
	a.NotNil(json.Unmarshal([]byte(`{"channel": {}}`), &r),
		"a non-string and non-number must produce an error.")
	a.EqualValues(0, r.Channel.Val)
	// test null.
	a.EqualValues(0, r.Nil.Val)
	a.EqualValues("0", r.Nil.Txt)

	val1 := unifi.NewFlexInt(5)
	val2 := unifi.NewFlexInt(4)
	val1.Add(val2)

	a.EqualValues(float64(9.0), val1.Val)
	a.EqualValues("9", val1.Txt)

	val1.AddFloat64(-4)
	a.EqualValues(float64(5.0), val1.Val)
	a.EqualValues("5", val1.Txt)

	val3 := *unifi.NewFlexInt(3)
	val3.AddFloat64(7)
	a.EqualValues(10, val3.Val)

	a.EqualValues(10, val3.Int())
	a.EqualValues(10, val3.Int64())
}

func TestFlexString(t *testing.T) {
	t.Parallel()
	a := assert.New(t)

	var r struct {
		JustString        unifi.FlexString `json:"just_string"`
		SingleArrayString unifi.FlexString `json:"single_array_string"`
		MultiArrayString  unifi.FlexString `json:"multi_array_string"`
	}
	// no errors unmarshalling
	a.Nil(json.Unmarshal([]byte(`{"just_string": "foo", "single_array_string": ["bar"], "multi_array_string": ["baz", "foo"]}`), &r))

	// simple string
	a.Equal("foo", r.JustString.Val)
	a.EqualValues([]string{"foo"}, r.JustString.Arr)

	// single array string
	a.Equal("bar", r.SingleArrayString.Val)
	a.EqualValues([]string{"bar"}, r.SingleArrayString.Arr)

	// multi array string
	a.Equal("baz, foo", r.MultiArrayString.Val)
	a.EqualValues([]string{"baz", "foo"}, r.MultiArrayString.Arr)
}
