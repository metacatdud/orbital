package components

import (
	"orbital/orbital"
)

func (comp *LoginComponent) BindEvents() {
	comp.di.Events().On("evt:auth:login:success", comp.eventLoginSuccess)
	comp.di.Events().On("evt:auth:login:fail", comp.eventLoginFail)
}

func (comp *LoginComponent) UnbindEvents() {
	comp.di.Events().Remove("evt:auth:login:success")
	comp.di.Events().Remove("evt:auth:login:fail")
}

func (comp *LoginComponent) eventLoginSuccess() {
	comp.di.State().Set("state:orbital:authenticated", true)
}

func (comp *LoginComponent) eventLoginFail(errRes *orbital.ErrorResponse) {
	comp.di.State().Set("state:auth:errored", ErrorManagerFields{
		Type:    errRes.Type,
		Message: errRes.Msg,
	})
}
