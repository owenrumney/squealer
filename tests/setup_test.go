package tests

import (
	"fmt"
	"os"
	"testing"
)

var testPath string

func TestMain(t *testing.M) {
	repo, err := unpackTestRepo("../test_resources/sloppygit.tar")
	defer func() { _ = os.RemoveAll(repo) }()
	if err != nil {
		panic(err)
	}
	testPath = fmt.Sprintf("%s/sloppy", repo)

	os.Exit(t.Run())
}
