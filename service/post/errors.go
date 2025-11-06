// service/post/errors.go
package postsvc

import "errors"

var (
	ErrBadInput = errors.New("bad input")
	ErrNotOwner = errors.New("not owner")
	ErrNotFound = errors.New("post not found")
)
