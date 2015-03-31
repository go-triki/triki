package log

// triki codes (TCs), for infos <0
const (
	okLoginTC  trikiCode = -100
	okSignupTC trikiCode = -150
)

// LoginOK informs the log about successful login.
func LoginOK(v ...interface{}) *Error {
	return &Error{
		What:      "user logged in",
		TrikiCode: okLoginTC,
		Info:      v,
	}
}

// SignupOK informs the log about successful sign-up of a user.
func SignupOK(v ...interface{}) *Error {
	return &Error{
		What:      "user signed-up",
		TrikiCode: okSignupTC,
		Info:      v,
	}
}
