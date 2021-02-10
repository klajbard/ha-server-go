package consumption

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func GetAllCons(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	device := ps.ByName("device")

	cons, err := AllCons(device)
	if err != nil {
		http.Error(w, http.StatusText(500)+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cons)

}

func GetCons(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	device := ps.ByName("device")
	date := ps.ByName("date")

	cons := OneCons(device, date)
	data := r.URL.Query()["data"][0]

	if data == "watt" {
		io.WriteString(w, strconv.FormatFloat(cons.Watt, 'f', 1, 64))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cons)
}

func PutCons(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cons, err := UpdateCons(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(406), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cons)
}
