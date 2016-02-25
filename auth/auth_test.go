package auth

import (
	"os/user"
	"testing"
)

//
func TestPostgresql(test *testing.T) {
	usr, err := user.Current()
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	pg, err := NewPostgresql(usr.Username, "postgres", "127.0.0.1:5432")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	if err = pg.Clear(); err != nil {
		test.Log(err)
		test.FailNow()
	}
	testDb(test, pg)
}

//
func TestMemory(test *testing.T) {
	memory := NewMemory()
	testDb(test, memory)
}

//
func testDb(test *testing.T, auth Authenticator) {

	tags, err := auth.GetTagsForToken("token")
	if err == nil {
		test.Log("there should have been an error")
		test.FailNow()
	}
	if len(tags) != 0 {
		test.Log("wrong number of tags were returned")
		test.FailNow()
	}

	err = auth.AddToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	err = auth.AddTags("token", []string{"a", "b"})
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	tags, err = auth.GetTagsForToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
	if len(tags) != 2 {
		test.Log("wrong number of tags were returned", tags)
		test.FailNow()
	}

	err = auth.RemoveTags("token", []string{"a", "b"})
	if err != nil {
		test.Log(err)
		test.FailNow()
	}

	err = auth.RemoveToken("token")
	if err != nil {
		test.Log(err)
		test.FailNow()
	}
}
