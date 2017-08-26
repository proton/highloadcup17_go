package main

import (
	"io"
	"strconv"
)

type Entity interface {
	// Update(data *JsonData, lock bool)
	UpdateFromJSON(data []byte, lock bool)
	writeJSON(w io.Writer)
}

type EntityRepo interface {
	// Create(data *JsonData)
	CreateFromJSON(data []byte)
	FindEntity(id int) (Entity, bool)
}

func find_entity(repo EntityRepo, entity_id_str *string) (Entity, bool) {
	entity_id_int, error := strconv.Atoi(*entity_id_str)
	if error == nil {
		entity_id := int(entity_id_int)
		return repo.FindEntity(entity_id)
	}
	return nil, false
}
