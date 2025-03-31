package b64

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type encoder struct {
	encoding *base64.Encoding
	name     string
}

var encoders = []encoder{
	{encoding: base64.StdEncoding, name: "StdEncoding"},
	{encoding: base64.RawStdEncoding, name: "RawStdEncoding"},
	{encoding: base64.URLEncoding, name: "URLEncoding"},
	{encoding: base64.RawURLEncoding, name: "RawURLEncoding"},
}

type encodeDecodeError struct {
	name    string
	message string
}

func (e encodeDecodeError) Error() string {
	return fmt.Sprintf("%s error: %s", e.name, e.message)
}

func TryDecodeString(s string) (bs []byte, ok bool, err error) {
	bs = make([]byte, 0)
	ok = false
	errs := make([]encodeDecodeError, 0)

	for _, enc := range encoders {
		var err0 error
		bs, err0 = enc.encoding.DecodeString(s)
		if err0 != nil {
			errs = append(errs, encodeDecodeError{name: enc.name, message: err0.Error()})
		} else {
			ok = true
			break
		}
	}

	if len(errs) > 0 {
		msgs := make([]string, len(errs))
		for i := range errs {
			msgs[i] = errs[i].Error()
		}
		err = errors.New(strings.Join(msgs, " | "))
	}

	return bs, ok, err
}
