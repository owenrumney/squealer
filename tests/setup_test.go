package tests

import (
	"fmt"
	"os"
	"testing"
)

var gitTestPath string
var dirTestPath string

func TestMain(t *testing.M) {
	getwd, _ := os.Getwd()
	repo, err := unpackTestRepo(fmt.Sprintf("%s/../test_resources/sloppygit.tar", getwd))
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(repo) }()
	gitTestPath = fmt.Sprintf("%s/sloppy", repo)

	dir, err := unpackTestRepo(fmt.Sprintf("%s/../test_resources/sloppy.tar", getwd))
	if err != nil {
		panic(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()
	dirTestPath = fmt.Sprintf("%s/sloppy2", dir)

	os.Exit(t.Run())
}
