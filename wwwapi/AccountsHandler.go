package wwwapi

import (
	"bitbucket.org/kornel661/triki/gotriki/db"
	"bitbucket.org/kornel661/triki/gotriki/log"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

// AccountsIDGetHandler serves account information on /api/accounts/:account_id.
func AccountsIDGetHandler(w http.ResponseWriter, r *http.Request) {
	// log info
	var info bytes.Buffer
	defer func() { log.Infoln(info.String()) }()
	apiAccessLog(&info, r)

	account_id := mux.Vars(r)["account_id"]
	if !bson.IsObjectIdHex(account_id) {
		http.Error(w, "Invalid account ID.", http.StatusBadRequest)
		fmt.Fprintf(&info, " Invalid account ID: `%s`.", account_id)
		return
	}
	accID := bson.ObjectIdHex(account_id)

	// check permissions
	if !hasPermission(accID, r) {
		http.Error(w, "Access denied. You don't have permissions to access this account's information. Try logging in.", http.StatusForbidden)
		_, usrID := authenticatedUser(r)
		fmt.Fprintf(&info, " Access denied for %s.", usrID)
		return
	}

	// fetch user (account) information
	usr, err := db.UserFindByID(accID)
	if err != nil {
		http.Error(w, "Error finding account: "+err.Error()+".", http.StatusNotFound)
		fmt.Fprintf(&info, " Error finding account: %s.", err)
		return
	}

	// reply
	writeJSON(w, &info, usr)
}
