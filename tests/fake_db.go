package main

import (
    "os"
    "time"

    "../api"
    "github.com/Coff3e/Api"
)

func main() {
    con_str := "fake_db.db"
    router := church.Church{}
    _, err := router.SignDB(con_str, api.Sqlite)
    if (err != nil) {
        os.Exit(1)
    }

    api.Log("Database connected with sucess")

    joao := church.User{}
    joao.ModelType = "User"
    joao.Name = "Joao da Silva"
    joao.Email = "joao@j.com"
    joao.Phone = "99 8452 1107"
    joao.Genre = "M"
    joao.State = "married"
    joao.Born = time.Now()
    joao.SetPass("maria")

    maria := church.User{}
    maria.Name = "Maria da Silva"
    maria.Email = "maria@j.com"
    maria.Genre = "F"
    maria.Phone = "99 8452 1108"
    maria.State = "married"
    maria.Born = time.Now()
    maria.SetPass("joao")

    person := church.Role{}
    person.Name = "Person"

    poor := church.Role{}
    poor.Name = "Poor"

    maria.Create()
    joao.Create()

    poor.Create()
    person.Create()

    person.Sign(joao)
    poor.Sign(joao)

    person.Sign(maria)

}
