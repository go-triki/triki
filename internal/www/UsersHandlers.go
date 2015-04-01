package www

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"gopkg.in/triki.v0/internal/log"
	"gopkg.in/triki.v0/internal/models/user"
)

// UsersIDGet serves account information on /api/users/:user_id.
func UsersIDGet(cx context.Context, w http.ResponseWriter, r *http.Request) {
	accID, err := readID(mux.Vars(r)["user_id"])
	if err != nil {
		writeError(cx, w, err)
		return
	}

	// TODO check permissions

	// fetch user (account) information
	usr, err := user.GetByID(cx, accID)
	if err != nil {
		writeError(cx, w, err)
		return
	}

	// TODO strip information (e.g. if usr != logged-in usr)

	// reply
	err = writeJSON(cx, w, Resp{
		Users: []*user.T{usr},
	})
	if err != nil {
		// error writing response, just log
		log.Set(cx, err)
		return
	}
	// TODO log?
}
