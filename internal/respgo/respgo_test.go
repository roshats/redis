package respgo_test

import (
	"bufio"
	"errors"
	"io"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/roshats/redis/internal/respgo"
)

func TestEncode(t *testing.T) {
	cases := []struct {
		value []byte
		want  string
	}{
		{respgo.EncodeString("OK"), "+OK\r\n"},
		{respgo.EncodeString("中文"), "+中文\r\n"},
		{respgo.EncodeString(""), "+\r\n"},

		{respgo.EncodeError("Error message"), "-Error message\r\n"},

		{respgo.EncodeInt(1000), ":1000\r\n"},
		{respgo.EncodeInt(1456061893587000000), ":1456061893587000000\r\n"},
		{respgo.EncodeInt(-1), ":-1\r\n"},

		{respgo.EncodeBulkString("foobar"), "$6\r\nfoobar\r\n"},
		{respgo.EncodeBulkString("中文"), "$6\r\n中文\r\n"},
		{respgo.EncodeBulkString(""), "$0\r\n\r\n"},
		{respgo.EncodeNull(), "$-1\r\n"},

		{respgo.EncodeArray([][]byte{
			respgo.EncodeBulkString("foo"),
			respgo.EncodeBulkString("bar"),
		}), "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n"},

		{respgo.EncodeArray([][]byte{
			respgo.EncodeInt(1),
			respgo.EncodeInt(2),
			respgo.EncodeInt(3),
		}), "*3\r\n:1\r\n:2\r\n:3\r\n"},

		{respgo.EncodeArray([][]byte{
			respgo.EncodeInt(1),
			respgo.EncodeInt(2),
			respgo.EncodeInt(3),
			respgo.EncodeInt(4),
			respgo.EncodeBulkString("foobar"),
		}), "*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n"},

		{respgo.EncodeArray([][]byte{
			respgo.EncodeArray([][]byte{
				respgo.EncodeInt(1),
				respgo.EncodeInt(2),
				respgo.EncodeInt(3),
			}),
			respgo.EncodeArray([][]byte{
				respgo.EncodeString("Foo"),
				respgo.EncodeError("Bar"),
			}),
		}), "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n"},

		{respgo.EncodeArray([][]byte{
			respgo.EncodeBulkString("foo"),
			respgo.EncodeNull(),
			respgo.EncodeBulkString("bar"),
		}), "*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n"},

		{respgo.EncodeArray([][]byte{}), "*0\r\n"},

		{respgo.EncodeNullArray(), "*-1\r\n"},
	}

	for _, item := range cases {
		if string(item.value) != item.want {
			t.Errorf("Encode get: %q ,want: %q", item.value, item.want)
		}
	}
}

func TestEncodeStringPanic(t *testing.T) {
	defer func() {
		if r, want := recover(), "SimpleString cannot contain a CR or LF character"; r != want {
			t.Errorf("Encode get: %v ,want: %v", r, want)
		}
	}()
	respgo.EncodeString("Hello World\r")
}

func TestEncodeBulkStringPanic(t *testing.T) {
	defer func() {
		if r, want := recover(), "BulkString is over 512 MB"; r != want {
			t.Errorf("Encode get: %v ,want: %v", r, want)
		}
	}()
	b := make([]byte, 536870913)
	respgo.EncodeBulkString(string(b))
}

func equal(v1 interface{}, v2 interface{}) bool {
	switch v1 := v1.(type) {
	case error:
		v2, ok := v2.(error)
		if !ok {
			return false
		}
		return v1.Error() == v2.Error()
	case []interface{}:
		v2, ok := v2.([]interface{})
		if !ok || len(v1) != len(v2) {
			return false
		}
		for i, _ := range v1 {
			if !equal(v1[i], v2[i]) {
				return false
			}
		}
		return true
	default:
		return v1 == v2
	}
}

func TestDecodeRawString(t *testing.T) {
	value := "get foo\r\n"
	want := []interface{}{"get", "foo"}
	bufReader := bufio.NewReader(strings.NewReader(value))
	result, err := respgo.Decode(bufReader)
	if err != nil {
		t.Errorf("case %q get error %v", value, err)
		return
	}
	if !equal(result, want) {
		t.Errorf("case %q get %v want %v", value, result, want)
	}
}

func TestDecode(t *testing.T) {
	cases := []struct {
		value string
		want  interface{}
	}{
		{"+OK\r\n", "OK"},
		{"+中文\r\n", "中文"},
		{"+\r\n", ""},

		{"-Error message\r\n", errors.New("Error message")},

		{":1000\r\n", int64(1000)},
		{":1456061893587000000\r\n", int64(1456061893587000000)},
		{":-1\r\n", int64(-1)},

		{"$6\r\nfoobar\r\n", "foobar"},
		{"$6\r\n中文\r\n", "中文"},
		{"$0\r\n\r\n", ""},
		{"$-1\r\n", nil},

		{"*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			[]interface{}{"foo", "bar"}},
		{"*3\r\n:1\r\n:2\r\n:3\r\n",
			[]interface{}{int64(1), int64(2), int64(3)}},
		{"*5\r\n:1\r\n:2\r\n:3\r\n:4\r\n$6\r\nfoobar\r\n",
			[]interface{}{int64(1), int64(2), int64(3), int64(4), "foobar"}},
		{"*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n",
			[]interface{}{[]interface{}{int64(1), int64(2), int64(3)}, []interface{}{"Foo", errors.New("Bar")}}},
		{"*3\r\n$3\r\nfoo\r\n$-1\r\n$3\r\nbar\r\n",
			[]interface{}{"foo", nil, "bar"}},
		{"*0\r\n", []interface{}{}},
		{"*-1\r\n", nil},

		//{"get foo", []interface{}{"get", "foo"}},
	}

	// Single Decode
	for _, item := range cases {
		result, err := respgo.Decode(bufio.NewReader(strings.NewReader(item.value)))
		if err != nil {
			t.Errorf("case %q get error %v", item.value, err)
			continue
		}
		if !equal(result, item.want) {
			t.Errorf("case %q get %v want %v", item.value, result, item.want)
		}
	}

	// Multiple Decode
	multiCase := ""
	for _, item := range cases {
		multiCase += item.value
	}
	bufReader := bufio.NewReader(strings.NewReader(multiCase))
	for i := 0; i < len(cases); i++ {
		result, err := respgo.Decode(bufReader)
		item := cases[i]
		if err != nil {
			t.Errorf("case %q get error %v", item.value, err)
			continue
		}
		if !equal(result, item.want) {
			t.Errorf("case %q get %v want %v", item.value, result, item.want)
		}
	}

	// Chaos Decode
	repeats := 10000
	bufReader = bufio.NewReader(&FakeNetIO{buf: []byte(strings.Repeat(multiCase, repeats))})
	for j := 0; j < repeats; j++ {
		for i := 0; i < len(cases); i++ {
			result, err := respgo.Decode(bufReader)
			item := cases[i]
			if err != nil {
				t.Errorf("case %q get error %v", item.value, err)
				continue
			}
			if !equal(result, item.want) {
				t.Errorf("case %q get %v want %v", item.value, result, item.want)
			}
		}
	}
}

type FakeNetIO struct {
	offset int
	buf    []byte
}

func (r *FakeNetIO) Read(p []byte) (n int, err error) {
	if r.offset < len(r.buf) {
		n = rand.Intn(len(p))
		if n == 0 {
			n = len(p)
		}
		if r.offset+n > len(r.buf) {
			n = len(r.buf) - r.offset
		}
		time.Sleep(time.Duration(n) * time.Microsecond)
		copy(p, r.buf[r.offset:r.offset+n])
		r.offset += n
		return n, nil
	}
	return 0, io.EOF
}

func TestDecodeError(t *testing.T) {
	cases := [][2]string{
		{"", "EOF"},
		{"+\n", `line is too short: "+\n"`},
		{"+x\n", `invalid CRLF: "+x\n"`},
		{"!0\r\n", `invalid RESP type: "!"`},
		{":x\r\n", "strconv.ParseInt: parsing \"x\": invalid syntax"},
		{"$x\r\nfoobar\r\n", "strconv.Atoi: parsing \"x\": invalid syntax"},
		{"$536870913\r\nxx\r\n", `invalid Bulk Strings length: 536870913`},
		{"$6\r\nfoobarxx", `invalid CRLF: "foobarxx"`},
		{"$-2\r\nxx\r\n", `invalid Bulk Strings length: -2`},
		{"$6\r\nfoo\r\n", "unexpected EOF"},
		{"*x\r\n:1\r\n:2\r\n", "strconv.Atoi: parsing \"x\": invalid syntax"},
		{"*2\r\n:1\r\n", "EOF"},
	}
	for _, item := range cases {
		_, err := respgo.Decode(bufio.NewReader(strings.NewReader(item[0])))
		if err == nil {

			t.Errorf("respgo.Decode %q should return error", item[0])
		} else if err.Error() != item[1] {
			t.Errorf("respgo.Decode %q get: %q ,want: %q", item[0], err, item[1])
		}
	}
}
