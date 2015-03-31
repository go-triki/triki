package www

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// UsersIDGet serves account information on /api/accounts/:account_id.
func UsersIDGet(cx context.Context, w http.ResponseWriter, r *http.Request) {
	accID, err := readID(mux.Vars(r)["user_id"])
	if err != nil {
		log.Set(cx, err)
		// TODO
		return
	}

	// TODO check permissions

	// fetch user (account) information
	usr, err := user.GetByID(cx, accID)
	if err != nil {
		log.Set(cx, err)
		// TODO
		return
	}

	// reply
	err = writeJSON(cx, w, usr)
	if err != nil {
		log.Set(cx, err)
		return
	}
	// TODO log?
}
