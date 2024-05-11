package web

import (
	"github.com/gorilla/mux"
	"mux-crud/happiness"
)

//var testHandler = handlers.AuthHandler{}
//
//var testActions = happiness.CustomActions{
//
//	happiness.Action{
//		Name:    "/signup",
//		Handler: happiness.Depend(authHandler.SignupHandler),
//		Method:  happiness.POST,
//	},
//	happiness.Action{
//		Name:    "/login",
//		Handler: happiness.Depend(authHandler.LoginHandler),
//		Method:  happiness.POST,
//	},
//}

type (
	UserMeta struct {
		happiness.BaseModel `bson:",inline"`
		Name                string `json:"name" bson:"name"`
		FirstName           string `json:"first_name" bson:"first_name"`
		LastName            string `json:"last_name" bson:"last_name"`
		Email               string `json:"email" bson:"email"`
		Phone               string `json:"phone" bson:"phone"`
	}

	UserMetaSchema struct {
		Name      string `json:"name" bson:"name"`
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Email     string `json:"email" bson:"email"`
		Phone     string `json:"phone" bson:"phone"`
	}
)

func (s UserMetaSchema) Validate() error {
	if s.Name == "" {
		//return fmt.Errorf("name is required")

	}

	//validate := validator.New()
	//err := validate.Struct(s)
	//if err != nil {
	//	return happiness.NewValidationError(err.(validator.ValidationErrors))
	//}
	return nil
}

func (u UserMeta) GetModelName() string {
	return "user_meta"
}

var userMetaRepository = happiness.NewBaseRepository(UserMeta{})

func TestRoute(r *mux.Router) {
	//happiness.Crud(r, services.UserRepository, "/auth", authActions,
	//	happiness.WithoutCreate, happiness.WithoutList, happiness.WithoutFetch,
	//	happiness.WithoutUpdate)

	happiness.Crud(r, userMetaRepository, "/user_model", happiness.WithoutCreate,
		happiness.CreateDto(UserMetaSchema{}))
}
