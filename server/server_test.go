package server

import (
	"testing"
)

var (
	testToken = "token"
)

// TestStart tests the auth start process
func TestStart(t *testing.T) {

	// test for error if an auth is provided w/o a token
	// if err := Start([]string{"tcp://127.0.0.1:1445"}, ""); err == nil {
	// 	t.Fatalf("Expecting error!")
	// }

	//
	// if err := Start([]string{"tcp://127.0.0.1:1445"}, testToken); err != nil {
	// 	t.Fatalf("Unexpected error!")
	// }

	// test for error if authtoken does not match the token the server started with
	// if authtoken != testToken {
	// 	t.Fatalf("Unexpected token!")
	// }
}

//
// func TestAuthCommands(t *testing.T) {
//
// 	//
// 	commands := GenerateAuthCommands(auth.NewMemory())
//
// 	// the client parameter is not used.
// 	if res := commands["register"].Handle(nil, []string{"token", "1,2,3,4"}); testForError(res) {
// 		t.Fatalf(res)
// 	}
//
// 	//
// 	if res := commands["set"].Handle(nil, []string{"token", "a,b,c,d"}); testForError(res) {
// 		t.Fatalf(res)
// 	}
//
// 	//
// 	if tags := commands["tags"].Handle(nil, []string{"token"}); tags == "" {
// 		t.Fatalf("wrong tags were returned")
// 	}
//
// 	//
// 	if res := commands["unset"].Handle(nil, []string{"token", "a,b,c,d"}); testForError(res) {
// 		t.Fatalf(res)
// 	}
//
// 	//
// 	if res := commands["unregister"].Handle(nil, []string{"token"}); testForError(res) {
// 		t.Fatalf(res)
// 	}
//
// 	//
// 	if tags := commands["tags"].Handle(nil, []string{"token"}); !testForError(tags) {
// 		t.Fatalf("wrong number of tags returned", tags)
// 	}
// }

//
// func testForError(s string) bool {
// 	return strings.HasPrefix(s, "Error:")
// }
