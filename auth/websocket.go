package auth

import (
	"fmt"
	"net/http"
)

type (
	websocket struct{}
)

//
func NewWS() websocket {
	return websocket{}
}

//
func (ws websocket) AddToken(token string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) RemoveToken(token string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) AddTags(token string, tags []string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) RemoveTags(token string, tags []string) error {
	return fmt.Errorf("I dont do anything...")
}

//
func (ws websocket) GetTagsForToken(token string) ([]string, error) {
	return []string{}, fmt.Errorf("I dont do anything...")
}

//
func AuthenticateWebsocket() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request){
		
	}

	// fmt.Println("AUTH WEBSOCKET???", DefaultAuth)

	// return func(w http.ResponseWriter, r *http.Request) {
	//
	// 	fmt.Println("HERE!!?!?!", r.FormValue("x-auth-token"), r.Header.Get("x-auth-token"))
	//
	// 	//
	// 	var token string
	// 	switch {
	// 	case r.Header.Get("x-auth-token") != "":
	// 		token = r.Header.Get("x-auth-token")
	// 	case r.FormValue("x-auth-token") != "":
	// 		token = r.FormValue("x-auth-token")
	// 	default:
	// 		token = "unauthorized"
	// 	}
	//
	// 	fmt.Println("TOKEN??", token)
	//
	// 	// if they have no tags registered for the token, then they are not authorized
	// 	// to connect to mist
	// 	if tags, err := DefaultAuth.GetTagsForToken(token); err != nil || len(tags) == 0 {
	// 		fmt.Println("BRONK??", err)
	// 		w.WriteHeader(401)
	// 		return
	// 	}
	//
	// 	fmt.Println("HERE!")
	//
	// 	// overwrite the subscribe command so that we can add authentication to it.
	// 	mixins := map[string]Handler{
	// 		"subscribe": {0, handleWSAuthSubscribe(token, a)},
	// 	}
	//
	// 	//
	// 	wsUpgrade := ListenWS(mixins)
	// 	wsUpgrade(w, r)
	// }
}
