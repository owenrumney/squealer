package tests

import (
	"fmt"
	"os"
	"testing"
)

var testPath string

func TestMain(t *testing.M) {
	getwd, _ := os.Getwd()
	repo, err := unpackTestRepo(fmt.Sprintf("%s/../test_resources/sloppygit.tar", getwd))
	defer func() { _ = os.RemoveAll(repo) }()
	if err != nil {
		panic(err)
	}
	testPath = fmt.Sprintf("%s/sloppy", repo)

	os.Exit(t.Run())
}
