package userserver

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/klaital/comics/pkg/config"
	"net/http"
)

func New(cfg *config.Config) *restful.WebService {
	// Users API
	usersWS := restful.WebService{}
	usersWS.Path(cfg.BasePath + "/users").ApiVersion("1.0.0").Doc("CRUD API for User management")
	//usersWS.Route(usersWS.GET("/{userId}").
	//	//Filter(filters.RequireValidJWT).
	//	To(DescribeUserHandler).
	//	Doc("Fetch player profile").
	//	Param(usersWS.PathParameter("userId", "User ID from the database")).
	//	Produces(restful.MIME_JSON).
	//	Writes(DescribeUserResponse{}).
	//	Returns(http.StatusOK, "Fetched user profile", DescribeUserResponse{}))

	// TODO: document failure responses
	usersWS.Route(usersWS.POST("/").
		To(RegisterUserHandler(cfg)).
		Doc("New User Signup").
		Consumes(restful.MIME_JSON).
		Reads(RegisterUserRequest{}).
		Writes(RegisterUserResponse{}).
		Returns(http.StatusCreated, "User registered successfully", RegisterUserResponse{}))
	//usersWS.Route(usersWS.POST("/login").
	//	To(LoginUserHandler).
	//	Doc("Login User to get a JWT").
	//	Consumes(restful.MIME_JSON). // TODO: maybe directly accept a form POST?
	//	Reads(LoginUserRequest).
	//	Produces(restful.MIME_JSON).
	//	Writes(LoginUserResponse).
	//	Returns(http.StatusOK, "User Logged in successfully", LoginUserResponse{}))

	return &usersWS
}
