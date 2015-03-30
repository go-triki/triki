package www

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/triki.v0/internal/log"
)

// AccountsIDGetHandler serves account information on /api/accounts/:account_id.
func AccountsIDGetHandler(w http.ResponseWriter, r *http.Request) {
	// log info
	var info bytes.Buffer
	defer func() { log.Infoln(info.String()) }()
	apiAccessLog(&info, r)

	account_id := mux.Vars(r)["account_id"]
	if !bson.IsObjectIdHex(account_id) {
		Error(w, http.StatusBadRequest, statusResourceInvalid,
			"Invalid account ID.")
		fmt.Fprintf(&info, " Invalid account ID: `%s`.", account_id)
		return
	}
	accID := bson.ObjectIdHex(account_id)

	// check permissions
	if !hasPermission(accID, r) {
		Error(w, http.StatusForbidden, statusUnauthorized,
			"Access denied. You don't have permissions to access this account's information. Try logging in.")
		_, usrID := authenticatedUser(r)
		fmt.Fprintf(&info, " Access denied for %s.", usrID)
		return
	}

	// fetch user (account) information
	usr, err := db.UserFindByID(accID)
	if err != nil {
		Error(w, http.StatusNotFound, statusError,
			"Error finding account: %s.", err)
		fmt.Fprintf(&info, " Error finding account: %s.", err)
		return
	}

	// reply
	writeJSON(w, &info, usr)
}
