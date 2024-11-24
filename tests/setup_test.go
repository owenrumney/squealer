package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/liamg/tml"
)

var gitTestPath string
var dirTestPath string

func TestMain(t *testing.M) {
	tml.DisableFormatting()
	repo, err := unpackTestRepo("../test_resources/sloppygit.tar")
	if err != nil {
		panic(err)
	}
	gitTestPath = fmt.Sprintf("%s/sloppy", repo)

	dir, err := unpackTestRepo("../test_resources/sloppy.tar")
	if err != nil {
		panic(err)
	}

	dirTestPath = fmt.Sprintf("%s/sloppy2", dir)

	result := t.Run()

	os.RemoveAll(repo)
	os.RemoveAll(dir)

	os.Exit(result)
}
