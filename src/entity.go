package main

import (
	"io"
)

type Entity interface {
	Update(data *JsonData)
	to_json(w io.Writer)
}

type EntityRepo interface {
	Create(data *JsonData)
	CreateFromJson(raw_data []byte) error
	Find(id uint32) (Entity, bool)
}
