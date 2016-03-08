package auth

import (
	// "fmt"
	"net/url"
	"testing"

	_ "github.com/lib/pq"
)

var (
	testToken = "token"
	testTag1  = "hello"
	testTag2  = "world"
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

	//
	if err := auth.AddTags(testToken, []string{testTag1, testTag2}); err != nil {
		t.Fatalf(err.Error())
	}

	//
	tags, err = auth.GetTagsForToken(testToken)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if len(tags) != 2 {
		t.Fatalf("Wrong number of tags. Expecting 2 received %v", len(tags))
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
