package controllers

import (
	"net/http"

	"github.com/Mauricio-3107/lenslocked-v2/views"
)

func StaticHandler(tpl views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, nil)
	}
}
