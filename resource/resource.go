package resource

import (
	"github.com/gorilla/context"
	"net/http"
	"encoding/json"
	"QuickAndDirtyAuthZ/database"
	"github.com/gorilla/mux"
)

type CommonLib struct {
	DB           *database.InMemory
	DatabaseName string
}

func (cl *CommonLib) Get(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	obj, err := cl.DB.GetFromDatabase(cl.DatabaseName, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(obj)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func (cl *CommonLib) Delete(rw http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	id := vars["id"]

	err := cl.DB.DeleteFromDatabase(cl.DatabaseName, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
}

func (cl *CommonLib) List(rw http.ResponseWriter, req *http.Request) {

	transactions, err := cl.DB.ListDatabase(cl.DatabaseName)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(transactions)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}

func (cl *CommonLib) ListWithDiscretion(rw http.ResponseWriter, req *http.Request) {

	// Getting Scope from context
	usernameFromContext := context.Get(req, "user")
	roleFromContext := context.Get(req, "role")
	println("Context Value: " + roleFromContext.(string))

	var err error
	var transactions []*database.Object

	if roleFromContext.(string) == "admin" {
		transactions, err = cl.DB.ListDatabase(cl.DatabaseName)

		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		transaction, err := cl.DB.GetFromDatabase(cl.DatabaseName, usernameFromContext.(string))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		transactions = append(transactions, transaction)
	}

	response, err := json.Marshal(transactions)

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)
}
