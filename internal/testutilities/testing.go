package testutilities

import (
	"testing"
	"unicode"
	"unicode/utf8"

	"github.com/google/go-cmp/cmp"
)

type Check func(got any, err error, t *testing.T, args ...any)

func GotSuccess(_ any, err error, t *testing.T, args ...any) {
	t.Helper()

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func GotErrorMessage(expectedMessage string) Check {
	return func(_ any, err error, t *testing.T, args ...any) {
		t.Helper()

		if err == nil {
			t.Errorf("expected error %q, got nil", expectedMessage)
			return
		}

		if err.Error() != expectedMessage {
			t.Errorf("expected error %q, got %q", expectedMessage, err)
		}
	}
}

func GotResult(want any, options ...cmp.Option) Check {
	return func(got any, err error, t *testing.T, args ...any) {
		t.Helper()

		if diff := cmp.Diff(want, got, options...); diff != "" {
			t.Errorf("expected %v, got %v", want, got)
		}
	}
}

func IgnoreUnexportedFields() cmp.Option {
	return cmp.FilterPath(func(path cmp.Path) bool {
		sf, ok := path.Index(-1).(cmp.StructField)
		if !ok {
			return false
		}
		r, _ := utf8.DecodeRuneInString(sf.Name())
		return !unicode.IsUpper(r)
	}, cmp.Ignore())
}
