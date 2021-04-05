package models

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"wingrent/database"
	"wingrent/database/postgres"
)

type PlaneModel struct {
	ID *int `json:"id"`
	Name *string `json:"name"`
}

func (pm *PlaneModel) WriteToDB() error {
	query := "INSERT INTO plane_model(name) VALUES($1) RETURNING id"
	err := database.DBSingleton.GetConnection().QueryRow(query, pm.Name).Scan(&pm.ID)
	if err != nil {
		return err
	}
	return nil
}

func RetrievePlaneModelByName(name string) (*PlaneModel, error) {
	query := "SELECT id, name FROM plane_model WHERE name = $1"
	var pm PlaneModel
	err := database.DBSingleton.GetConnection().QueryRow(query, name).Scan(&pm.ID, &pm.Name)
	if err == sql.ErrNoRows{
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &pm, err
}

type Plane struct {
	ID *int `json:"id"`
	TailNumber *string `json:"tail_number"`
	Model PlaneModel `json:"model"`
}

func RetrievePlane(id int) (*Plane, error) {
	query := "SELECT p.id, p.tail_number, m.id, m.name FROM planes p JOIN plane_model m ON p.plane_model_id=m.id WHERE p.id = $1"
	var plane Plane
	err := database.DBSingleton.GetConnection().QueryRow(query, id).Scan(
		&plane.ID, &plane.TailNumber, &plane.Model.ID, &plane.Model.Name)
	if err != nil {
		return nil, err
	}
	return &plane, err
}

func (p *Plane) UpdateIdFromDBByUniqueFields() error {
	query := "SELECT id FROM planes WHERE tail_number = $1"
	err := database.DBSingleton.GetConnection().QueryRow(query, p.TailNumber).Scan(&p.ID)
	if err == sql.ErrNoRows{
		return nil
	}

	return err
}

func (p Plane) WriteToDB() (*Plane, error) {
	// get or create PlaneModel
	pm, _ := RetrievePlaneModelByName(*p.Model.Name)
	if pm == nil {
		pm = &PlaneModel{
			Name: p.Model.Name,
		}
		err := pm.WriteToDB()
		if err != nil {
			return nil, err
		}
	}
	p.Model = *pm

	// create new plane
	query := "INSERT INTO planes(id, tail_number, plane_model_id) VALUES(default, $1, $2) RETURNING id"
	err := database.DBSingleton.GetConnection().QueryRow(query, p.TailNumber, p.Model.ID).Scan(&p.ID)
	if err.(*pq.Error).Code == postgres.DuplicateKey {
		err = p.UpdateIdFromDBByUniqueFields()
		if err != nil {
			return nil, err
		}
		return &p, &database.EntityAlreadyExists{}
	}
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (p Plane) Delete() error {
	query := "DELETE FROM planes WHERE id=$1"
	_, err := database.DBSingleton.GetConnection().Exec(query, p.ID)
	return err
}

func (p *Plane) Update(id int) error{
	if p.TailNumber != nil {
		query := "UPDATE planes SET tail_number=$1 WHERE id=$2"
		_, err := database.DBSingleton.GetConnection().Exec(query, p.TailNumber, p.ID)
		if err != nil {
			return err
		}
	}

	if p.Model.Name != nil {
		query := "UPDATE plane_model SET name=$1 WHERE id=$2"
		_, err := database.DBSingleton.GetConnection().Exec(query, p.Model.Name, p.Model.ID)
		if err != nil {
			return err
		}
	}

	newP, err := RetrievePlane(*p.ID)
	if newP == nil {
		return errors.New("Cannot retrieve updated entity")
	}
	*p = *newP
	return err
}

type PlanesList struct {
	Planes []Plane `json:"planes"`
}

func ListPlanes() (*PlanesList, error) {
	list := &PlanesList{}
	query := "SELECT p.id, p.tail_number, m.id, m.name FROM planes p JOIN plane_model m ON p.plane_model_id=m.id ORDER BY p.id"
	rows, err := database.DBSingleton.GetConnection().Query(query)
	if err != nil {
		return list, err
	}

	for rows.Next() {
		var plane Plane
		err := rows.Scan(&plane.ID, &plane.TailNumber, &plane.Model.ID, &plane.Model.Name)
		if err != nil {
			return &PlanesList{}, err
		}
		list.Planes = append(list.Planes, plane)
	}
	return list, nil
}
