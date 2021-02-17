package church

import (
    "github.com/Coff3e/Api"
)

type UserCelule struct {
    api.UserGroup
}

type Celule struct {
    api.Group
    Type       string `json:"type,omitempty" gorm:"index"`
    AddrId     uint   `json:",empty" gorm:"index"`
    Leader     uint   `json:",empty" gorm:"index"`
    CoLeader   uint   `json:",empty" gorm:"index"`
}

func (self UserCelule) TableName() string {
    return "user_groups"
}
func (self Celule) TableName() string {
    return "groups"
}

func (model *Celule) Create() {
    model.ModelType = "Celule"

    db.Create(model)

    e := db.First(model)
    if e.Error == nil {

        ID := model.ID
        ModelType := model.ModelType
        api.Log("Created", api.ToLabel(ID, ModelType))
    }
}

func (model *UserCelule) Create() {
    model.ModelType = "UserCelule"

    db.Create(model)

    e := db.First(model)
    if e.Error == nil {

        ID := model.ID
        ModelType := model.ModelType
        api.Log("Created", api.ToLabel(ID, ModelType))
    }
}


func (self UserCelule) Sign(user User, celule Celule) (User, Celule) {
    self.UserId = user.ID
    self.GroupId = celule.ID

    self.Create()
    api.Log("Linked", api.ToLabel(user.ID, user.ModelType), user.Name, "to", api.ToLabel(celule.ID, celule.ModelType), celule.Name)

    return user, celule
}

func (self UserCelule) Unsign(user User, celule Celule) (User, Celule) {
    self.UserId = user.ID
    self.GroupId = celule.ID

    self.Delete()
    api.Log("Unlinked", api.ToLabel(user.ID, user.ModelType), user.Name, "from", api.ToLabel(celule.ID, celule.ModelType), celule.Name)

    return user, celule
}

func (self Celule) Sign(user User) (User, Celule) {
    link := UserCelule{}
    user, self = link.Sign(user, self)

    return user, self
}

func (self Celule) Unsign(user User) (User, Celule) {
    link := UserCelule{}
    e := db.Where("user_id = ? AND group_id = ?", user.ID, self.ID).First(&link)
    if e.Error == nil {
        user, self = link.Unsign(user, self)
    }

    return user, self
}

func (self *Celule) GetUsers() []User {
    e := db.First(self)
    if e.Error == nil {
        user_list := []uint{}
        users := []User{}
        e := db.Raw("SELECT u.id FROM users u INNER JOIN user_groups ur INNER JOIN groups r ON ur.group_id = r.id AND ur.user_id = u.id AND r.id = ?", self.ID).Find(&user_list)
        if e.Error == nil {
            e := db.Find(&users, "id in ?", user_list)
            if e.Error == nil {
                return users
            }
        }
    }

    return []User{}
}

func (self *User) GetCelules() []Celule {
    e := db.First(self)
    if e.Error == nil {
        celule_list := []uint{}
        celules := []Celule{}
        e := db.Raw("SELECT r.id FROM groups r INNER JOIN user_groups ur INNER JOIN users u ON ur.group_id = r.id AND ur.user_id = u.id AND u.id = ?", self.ID).Find(&celule_list)
        if e.Error == nil {
            e := db.Find(&celules, "id in ?", celule_list)
            if e.Error == nil {
                return celules
            }
        }
    }

    return []Celule{}
}
