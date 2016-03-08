package auth

import (
	"net/url"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testToken = "token"
	testTag1  = "onefish"
	testTag2  = "twofish"
	testTag3  = "redfish"
	testTag4  = "bluefish"
)

// TestStart
func TestStart(t *testing.T) {

	//
	if err := Start("memory://", ""); err == nil {
		t.Fatalf("Expecting error!")
	}

	//
	if err := Start("memory://", "TOKEN"); err != nil {
		t.Fatalf(err.Error())
	}

	//
	if DefaultAuth == nil {
		t.Fatalf("Unexpected nil DefaultAuth!")
	}

	//
	if Token == "" {
		t.Fatalf("Unexpected blank Token!")
	}
}

// TestMemory
func TestMemory(t *testing.T) {

	//
	url, err := url.Parse("memory://")
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	mem, err := NewMemory(url)
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	testAuth(mem, t)
}

// TestScribble
func TestScribble(t *testing.T) {

	// attempt to remove the db from any previous tests
	if err := os.RemoveAll("/tmp/scribble"); err != nil {
		t.Fatalf(err.Error())
	}

	//
	url, err := url.Parse("scribble://?db=/tmp/scribble")
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	scribble, err := NewScribble(url)
	if err != nil {
		t.Fatalf(err.Error())
	}

	//
	testAuth(scribble, t)
}

// TestPostgres skip this for now because I don't want to have a postgres running
// func TestPostgres(t *testing.T) {
//
// 	//
// 	url, err := url.Parse("postgres://postgres@127.0.0.1:5432?db=postgres")
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	pg, err := NewPostgres(url)
// 	if err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	if _, err := pg.(postgresql).exec("TRUNCATE tokens, tags"); err != nil {
// 		t.Fatalf(err.Error())
// 	}
//
// 	//
// 	testAuth(pg, t)
// }

// testAuth
func testAuth(auth Authenticator, t *testing.T) {

	//
	tags, err := auth.GetTagsForToken(testToken)
	if err == nil {
		t.Fatalf("Expecting error!")
	}
	if len(tags) != 0 {
		t.Fatalf("Wrong number of tags. Expecting 0 got %v", len(tags))
	}

	//
	if err := auth.AddToken(testToken); err != nil {
		t.Fatalf(err.Error())
	}

	// add tags
	if err := auth.AddTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	// add same tags; these should not get added
	if err := auth.AddTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	// add same tags, different order; these should not get added
	if err := auth.AddTags(testToken, []string{testTag2, testTag1}); err != nil {
		t.Fatalf(err.Error())
	}

	// add more tags
	if err := auth.AddTags(testToken, []string{testTag3, testTag4}); err != nil {
		t.Fatalf(err.Error())
	}

	// this tests to ensure that same tags don't get added and we only get back
	// what we expect (in this case 4 unique tags)
	tags, err = auth.GetTagsForToken(testToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tags) != 4 {
		t.Fatalf("Wrong number of tags. Expecting 4 received %v", len(tags))
	}

	//
	if err := auth.RemoveTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	//
	if err := auth.RemoveToken(testToken); err != nil {
		t.Fatalf(err.Error())
	}
}
