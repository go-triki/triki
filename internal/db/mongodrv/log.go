package mongo

// Log saves record in the `log` collection in the DB.
func Log(record map[string]interface{}) error {
	return logSession.DB("").C(logCName).Insert(record)
}
