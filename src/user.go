package main

import (
	"encoding/json"
	"io"
	"sync"
)

type User struct {
	Id        uint32       `json:"id"`
	Email     string       `json:"email"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Gender    string       `json:"gender"`
	BirthDate int32        `json:"birth_date"`
	Mutex     sync.RWMutex `json:"-"`
	// VisitIds  []uint32
}

type UsersRepo struct {
	Collection map[uint32]*User
	Mutex      sync.RWMutex
}

func (entity *User) Update(data *JsonData) {
	entity.Mutex.Lock()
	for key, value := range *data {
		switch key {
		case "id":
			entity.Id = uint32(value.(float64))
		case "email":
			entity.Email = value.(string)
		case "first_name":
			entity.FirstName = value.(string)
		case "last_name":
			entity.LastName = value.(string)
		case "gender":
			entity.Gender = value.(string)
		case "birth_date":
			entity.BirthDate = int32(value.(float64))
			// bday := time.Unix(int64(entity.BirthDate), 0)
			// now := time.Now()
			// fmt.Println("Time now is:", now)
			// age := now.Year() - bday.Year()
			// if now.Month() < bday.Month() {
			//  age = age - 1
			// } else if (now.Month() == bday.Month()) && (now.Day() < bday.Day()) {
			//  age = age - 1
			// }
			// fmt.Println("User:", (*data)["id"])
			// fmt.Println("Age is:", age)

			// entity.Age = uint32(age)
		}
	}
	entity.Mutex.Unlock()
}

func (entity *User) to_json(w io.Writer) {
	entity.Mutex.RLock()
	json.NewEncoder(w).Encode(entity)
	entity.Mutex.RUnlock()
}

func NewUsersRepo() UsersRepo {
	return UsersRepo{
		Collection: make(map[uint32]*User),
		Mutex:      sync.RWMutex{}}
}

func (repo *UsersRepo) Create(data *JsonData) {
	entity := &User{}
	entity.Update(data)
	repo.Add(entity)
}

func (repo *UsersRepo) CreateFromJson(raw_data []byte) error {
	entity := &User{}
	err := json.Unmarshal(raw_data, entity)
	if err == nil {
		repo.Add(entity)
	}
	return err
}

func (repo *UsersRepo) Add(entity *User) {
	repo.Mutex.Lock()
	repo.Collection[entity.Id] = entity
	repo.Mutex.Unlock()
}

func (repo *UsersRepo) Find(id uint32) (Entity, bool) {
	repo.Mutex.RLock()
	var entity, found = repo.Collection[id]
	repo.Mutex.RUnlock()
	return entity, found
}
