package main

import (
	"io"
)

type Entity interface {
	Update(data *JsonData, lock bool) bool
	to_json(w io.Writer)
}

type EntityRepo interface {
	Create(data *JsonData) bool
	// CreateFromJson(raw_data []byte) error
	FindEntity(id uint32, lock bool) (Entity, bool)
}
