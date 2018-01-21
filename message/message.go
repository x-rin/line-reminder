package message

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func GetMessage (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := "Get Message!"
	json.NewEncoder(w).Encode(message)
}

func PostMessage (w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := "Posted Message: " + r.PostFormValue("body")
	json.NewEncoder(w).Encode(message)
}
