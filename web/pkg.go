/*
Package web is responsible for logging users in to the validator service,
listing the user's repositories and their validation status/results, enabling
and disabling hooks on the GIN server running the validation.
*/
package web

import gweb "github.com/G-Node/gin-cli/web"

const (
	serveralias = "gin"
)

var (
	sessions = make(map[string]*usersession)
	hookregs = make(map[string]gweb.UserToken)
)

type usersession struct {
	sessionID string
	gweb.UserToken
}
