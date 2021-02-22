package main

import (
    "os"

    "./api"
    "github.com/Coff3e/Api"
)

func main() {
    con_str := "host=localhost user=plankiton password=joaojoao dbname=church port=5432 sslmode=disable TimeZone=America/Araguaina"
    r := church.Church{}
    _, err := r.SignDB(con_str, api.Postgres)
    if (err != nil) {
        os.Exit(1)
    }
    api.Log("Database connected with sucess")
    r.
    Add(
        "post", "/login", api.RouteConf {
            "need-auth": false,
        }, church.LogIn,
    ).
    Add(
        "post", "/logout", nil, church.LogOut,
    ).
    Add(
        "post", "/verify", nil, church.Verify,
    ).
    Add(
        "post", "/user", nil, church.CreateUser,
    ).
    Add(
        "post", "/user/{id}/profile", nil, church.CreateUserProfile,
    )
    r.Run("/", 8000)
}
