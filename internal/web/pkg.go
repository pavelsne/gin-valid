/*
Package web is responsible for logging users in to the validator service,
listing the user's repositories and their validation status/results, enabling
and disabling hooks on the GIN server running the validation.
*/
package web

const (
	serveralias = "gin"
	progressmsg = "A validation job for this repository is currently in progress, please do not leave this page and refresh the page after a while."
	notfound = "This repository has not been validated yet. <a href='/pubvalidate' alt='onetime'>Back to one-time-validation</a>"
)
