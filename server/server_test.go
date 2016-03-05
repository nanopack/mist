package server

// import (
// 	"strings"
// 	"testing"
//
// 	"github.com/nanopack/mist/auth"
// )

//
// func TestAuthCommands(test *testing.T) {
//
// 	//
// 	commands := GenerateAuthCommands(auth.NewMemory())
//
// 	// the client parameter is not used.
// 	if res := commands["register"].Handle(nil, []string{"token", "1,2,3,4"}); testForError(res) {
// 		test.Log(res)
// 		test.FailNow()
// 	}
//
// 	//
// 	if res := commands["set"].Handle(nil, []string{"token", "a,b,c,d"}); testForError(res) {
// 		test.Log(res)
// 		test.FailNow()
// 	}
//
// 	//
// 	if tags := commands["tags"].Handle(nil, []string{"token"}); tags == "" {
// 		test.Log("wrong tags were returned")
// 		test.FailNow()
// 	}
//
// 	//
// 	if res := commands["unset"].Handle(nil, []string{"token", "a,b,c,d"}); testForError(res) {
// 		test.Log(res)
// 		test.FailNow()
// 	}
//
// 	//
// 	if res := commands["unregister"].Handle(nil, []string{"token"}); testForError(res) {
// 		test.Log(res)
// 		test.FailNow()
// 	}
//
// 	//
// 	if tags := commands["tags"].Handle(nil, []string{"token"}); !testForError(tags) {
// 		test.Log("wrong tags were returned", tags)
// 		test.FailNow()
// 	}
// }

//
// func testForError(s string) bool {
// 	return strings.HasPrefix(s, "Error:")
// }
