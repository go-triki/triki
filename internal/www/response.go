package www

import (
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/token"
	"gopkg.in/triki.v0/internal/models/user"
)

type Resp struct {
	Users  []*user.T    `json:"usrs,omitempty"`
	Tokens []*token.T   `json:"tkns,omitempty"`
	Errors []*log.Error `json:"errs,omitempty"`
}
