package util

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
)

const exitCodeFailure = 1

func Fail(msg any) {
	fmt.Fprintf(os.Stderr, "%v\n", msg)
	os.Exit(exitCodeFailure)
}

func FailOnError(err error) {
	if err != nil {
		Fail(err)
	}
}

func Must[T any](t T, err error) T {
	FailOnError(err)
	return t
}

func MustOK[T any](t T, ok bool) T {
	if !ok {
		Fail(fmt.Sprintf("MustOK: %#v failed", t))
	}
	return t
}

func NoError[T any](_ T, err error) bool { return err == nil }
func OK[T any](_ T, ok bool) bool        { return ok }

func Hash(v string) string { return fmt.Sprintf("%x", sha256.Sum256([]byte(v))) }

func TempDir(p ...string) string {
	return filepath.Join(append([]string{os.TempDir(), "ndql"}, p...)...)
}
