package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"wingrent/database"
	"wingrent/models"
)

func ListPlanes(w http.ResponseWriter, r *http.Request)  {
	planes, err := models.ListPlanes()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server error")
	}
	RespondWithJSON(w, http.StatusOK, planes)
}

func RetrievePlane(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	p, err := models.RetrievePlane(id)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	RespondWithJSON(w, http.StatusOK, p)
}

func CreatePlane(w http.ResponseWriter, r *http.Request) {
	p := &models.Plane{}
	err := json.NewDecoder(r.Body).Decode(p)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Incorrect payload")
		return
	}
	defer r.Body.Close()

	p, err = p.WriteToDB()
	if errors.Is(err, &database.EntityAlreadyExists{}) {
		RespondWithJSON(w, http.StatusOK, p)
		return
	}
	if err != nil{
		RespondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	RespondWithJSON(w, http.StatusCreated, p)
}


func DeletePlane(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid id")
		return
	}

	p := models.Plane{ID: &id}
	err = p.Delete()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	RespondWithJSON(w, http.StatusNoContent, map[string]string{"result": "success"})
}

func UpdatePlane(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var p models.Plane
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	p.ID = &id
	if err := p.Update(id); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Server Error")
		return
	}

	RespondWithJSON(w, http.StatusOK, p)
}
