package routes

import (
	"21-api/config"
	comment "21-api/features/comment"
	user "21-api/features/user"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func InitRoute(c *echo.Echo, ctl user.Controller, cc comment.CommentController) {
	userRoute(c, ctl)
	//activityRoute(c, ac)

}

func userRoute(c *echo.Echo, ctl user.Controller) {
	c.POST("/login", ctl.Login())

	//Register User
	c.POST("/register", ctl.Register()) //Endpoint untuk API

	//DeleteUser
	c.DELETE("/users/:id", ctl.Delete(), echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(config.JWTSECRET),
	}))

	// UpdateUser
	c.PUT("/users/:hp", ctl.Update())

}

func commentRoute(c *echo.Echo, cc comment.CommentController) {
	//Menambahkan Komentar
	c.POST("/comment", cc.AddComment())

	//Delete Komentar
	c.DELETE("/comment/:commentID", cc.DeleteComment())
}
