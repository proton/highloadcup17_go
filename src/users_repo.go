package main

import (
	"strconv"
	"sync"
)

type UsersRepo struct {
	Collection map[uint32]*User
	Mutex      sync.RWMutex
}

func (repo *UsersRepo) InitEntity() *User {
	return &User{}
}

func (repo *UsersRepo) CreateFromJSON(data []byte) {
	entity := repo.InitEntity()
	entity.UpdateFromJSON(data, false)
	repo.Add(entity)
}

func (repo *UsersRepo) Add(entity *User) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *UsersRepo) Find(id uint32) (*User, bool) {
	repo.Mutex.RLock()
	defer repo.Mutex.RUnlock()
	entity, ok := repo.Collection[id]
	return entity, ok
}

func (repo *UsersRepo) FindEntity(id uint32) (Entity, bool) {
	return repo.Find(id)
}

func find_user(entity_id_str []byte) (*User, bool) {
	entity_id_int, error := strconv.Atoi(bstring(entity_id_str))
	if error == nil {
		entity_id := uint32(entity_id_int)
		return Users.Find(entity_id)
	}
	return nil, false
}
