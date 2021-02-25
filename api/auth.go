package church

import (
    "github.com/Coff3e/Api"
    "fmt"
)

type Token struct {
    api.Token
}

func (model *Token) Create() {
    model.ModelType = api.GetModelType(model)

    user := User{}
    user.ID = model.UserId

    e := db.First(&user)
    if e.Error == nil {
        var order int64
        db.Find(model).Count(&order)

        model.ID = api.ToHash(fmt.Sprintf(
            "%d;%d;%s;%s;%s", order, user.ID, user.Name, user.Email, user.Phone,
        ))

        db.Create(model)
        e = db.First(model)
        if e.Error == nil {
            ID := model.ID
            ModelType := model.ModelType
            api.Log("Created", api.ToLabel(ID, ModelType))
        }
    }
}

func LogIn(r api.Request) (api.Response, int) {
    var data map[string]interface{}
    if validData(r.Data, generic_json_obj) {
        data = r.Data.(map[string]interface{})
    } else {
        msg := fmt.Sprint("Authentication fail, bad request, data need to be a object")
        api.Err(msg)
        return api.Response {
            Message: msg,
            Type:    "Error",
        }, 400
    }

    user := User {}

    res := db.First(&user, "email = ?", data["email"])
    if res.Error != nil {
        msg := fmt.Sprint("Authentication fail, no users with \"", data["email"], "\" registered")
        api.Err(msg)
        return api.Response {
            Message: msg,
            Type:    "Error",
        }, 404
    }

    if user.CheckPass(data["pass"].(string)) {
        token := Token {}
        token.UserId = user.ID
        token.Create()

        return api.Response {
            Type: "Sucess",
            Data: map[string]interface{} {
                "token": token.ID,
                "user": user,
            },
        }, 200
    }

    msg := fmt.Sprint("Authentication fail, password from <", user.Name, ":", user.ID,"> is wrong, permission denied")
    api.Err(msg)
    return api.Response {
        Message: msg,
        Type:    "Error",
    }, 405
}

func Verify(r api.Request) (api.Response, int) {
    token := Token {}
    token.ID = r.Token

    user := User {}
    user.ID = token.UserId


    db.First(&token)
    res := db.First(&user)
    if res.Error == nil {
        msg := fmt.Sprint("Token \"", r.Token, "\" Is valid")
        api.Log(msg)
        return api.Response {
            Type: "Sucess",
            Message: msg,
        }, 200
    }

    msg := fmt.Sprint("Authentication fail, user not found, permission denied")
    api.Err(msg)
    token.Delete()
    return api.Response {
        Message: msg,
        Type:    "Error",
    }, 404
}

func LogOut(r api.Request) (api.Response, int) {
    token := Token {}
    token.ID = r.Token

    db.First(&token)
    token.Delete()

    msg := fmt.Sprint("Token \"", r.Token, "\" removed")
    api.Log(msg)
    return api.Response {
        Type: "Sucess",
        Message: msg,
    }, 200
}
