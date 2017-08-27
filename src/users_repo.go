package main

import (
	"strconv"
	"sync"
)

var (
	USERS_REPO_COLLECTION_SIZE = uint32(1010000)
)

type UsersRepo struct {
	Collection []*User
	// Mutex         sync.RWMutex
	MapCollection map[uint32]*User
	MapMutex      sync.RWMutex
}

func (repo *UsersRepo) InitEntity() *User {
	entity := User{}
	return &entity
}

func (repo *UsersRepo) Create(data []byte) {
	entity := repo.InitEntity()
	entity.Update(data, false)
	repo.Add(entity)
}

func (repo *UsersRepo) Add(entity *User) {
	if entity.Id < USERS_REPO_COLLECTION_SIZE {
		repo.Collection[entity.Id] = entity
	} else {
		repo.MapMutex.Lock()
		defer repo.MapMutex.Unlock()
		repo.MapCollection[entity.Id] = entity
	}
}

func (repo *UsersRepo) Find(id uint32) (*User, bool) {
	if id < USERS_REPO_COLLECTION_SIZE {
		entity := repo.Collection[id]
		return entity, (entity != nil)
	} else {
		repo.MapMutex.Lock()
		defer repo.MapMutex.Unlock()
		entity, ok := repo.MapCollection[id]
		return entity, ok
	}
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
