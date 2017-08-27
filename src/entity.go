package main

import (
	"io"
	"strconv"
)

type Entity interface {
	UpdateFromJSON(data []byte, lock bool)
	writeJSON(w io.Writer)
}

type EntityRepo interface {
	CreateFromJSON(data []byte)
	FindEntity(id uint32) (Entity, bool)
}

func find_entity(repo EntityRepo, entity_id_str []byte) (Entity, bool) {
	entity_id_int, error := strconv.Atoi(bstring(entity_id_str))
	if error == nil {
		entity_id := uint32(entity_id_int)
		return repo.FindEntity(entity_id)
	}
	return nil, false
}
