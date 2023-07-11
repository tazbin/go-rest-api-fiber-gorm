package routes

import (
	"errors"
	"rest-api/database"
	"rest-api/model"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func CreateUserResponse(user *model.User) *User {
	return &User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
}

func CreateUser(c *fiber.Ctx) error {
	var user model.User

	if err := c.BodyParser(&user); err != nil {
		c.Status(400).JSON(err.Error())
	}

	database.Database.DB.Create(&user)
	createUserResponse := CreateUserResponse(&user)

	return c.Status(200).JSON(createUserResponse)
}

func GetAllUsers(c *fiber.Ctx) error {
	users := []model.User{}
	database.Database.DB.Find(&users)

	responseUsers := []*User{}

	for _, user := range users {
		responseUser := CreateUserResponse(&user)
		responseUsers = append(responseUsers, responseUser)
	}

	return c.Status(200).JSON(responseUsers)
}

func GetUserById(c *fiber.Ctx) error {
	var user model.User

	userId, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON("invalid user id")
	}

	database.Database.DB.Find(&user, userId)

	if user.ID == 0 {
		err := errors.New("no user found with this id")
		return c.Status(400).JSON(err.Error())
	}

	userResponse := CreateUserResponse(&user)
	return c.Status(200).JSON(userResponse)
}

func UpdateUser(c *fiber.Ctx) error {
	var user model.User

	userId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("invalid user id")
	}

	database.Database.DB.Find(&user, userId)

	if user.ID == 0 {
		return c.Status(400).JSON("no user found with this id")
	}

	type updateUserData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	var updateData updateUserData
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	user.FirstName = updateData.FirstName
	user.LastName = updateData.LastName

	database.Database.DB.Save(&user)
	userUpdateResponse := CreateUserResponse(&user)

	return c.Status(200).JSON(userUpdateResponse)
}

func DeleteUser(c *fiber.Ctx) error {
	var user model.User

	userId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON("invalid user id")
	}

	if err := database.Database.DB.Delete(&user, userId).Error; err != nil {
		return c.Status(400).JSON("user do not exist with that id")
	}

	return c.Status(200).JSON("user deleted successfully!")
}
