package routes

import (
	"fmt"
	"rest-api/database"
	"rest-api/model"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID        uint      `json:"id"`
	Product   Product   `json:"product"`
	User      User      `json:"user"`
	CreatedAt time.Time `json:"order_date"`
}

func createOrderResponse(order *model.Order, product *model.Product, user *model.User) Order {
	return Order{
		ID:        order.ID,
		Product:   *createProductResponse(product),
		User:      *CreateUserResponse(user),
		CreatedAt: order.CreatedAt,
	}
}

func CreateOrder(c *fiber.Ctx) error {
	order := &model.Order{}

	if err := c.BodyParser(order); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	product := &model.Product{}
	if err := database.Database.DB.Find(product, order.ProductRefer).Error; err != nil || product.ID == 0 {
		return c.Status(404).JSON(fmt.Sprintf("no product found with id %v", order.ProductRefer))
	}

	user := &model.User{}
	if err := database.Database.DB.Find(user, order.UserRefer).Error; err != nil || user.ID == 0 {
		return c.Status(404).JSON(fmt.Sprintf("no user found with id %v", order.UserRefer))
	}

	if err := database.Database.DB.Create(order).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	rOrder := createOrderResponse(order, product, user)
	return c.Status(200).JSON(rOrder)
}

func GetOrders(c *fiber.Ctx) error {
	orders := []*model.Order{}

	if err := database.Database.DB.Find(&orders).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	rAllOrders := []*Order{}
	for _, order := range orders {
		user := model.User{}
		if err := database.Database.DB.Find(&user, order.UserRefer).Error; err != nil {
			return c.Status(500).JSON(err.Error())
		}

		product := model.Product{}
		if err := database.Database.DB.Find(&product, order.ProductRefer).Error; err != nil {
			return c.Status(500).JSON(err.Error())
		}

		rOrder := createOrderResponse(order, &product, &user)
		rAllOrders = append(rAllOrders, &rOrder)
	}

	return c.Status(200).JSON(rAllOrders)
}

func GetOrderById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	order := &model.Order{}
	if err = database.Database.DB.First(order, id).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if order.ID == 0 {
		return c.Status(400).JSON("order not found with that id")
	}

	product := &model.Product{}
	if err = database.Database.DB.First(product, order.ProductRefer).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if product.ID == 0 {
		return c.Status(400).JSON("product not found with that id")
	}

	user := &model.User{}
	if err = database.Database.DB.First(user, order.UserRefer).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	if user.ID == 0 {
		return c.Status(400).JSON("user not found with that id")
	}

	rOrder := createOrderResponse(order, product, user)
	return c.Status(200).JSON(rOrder)
}

func UpdateOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type oBody struct {
		ProductId uint `json:"product_id"`
		UserId    uint `json:"user_id"`
	}

	orderBody := &oBody{}
	if err = c.BodyParser(orderBody); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	order := &model.Order{}
	if err = database.Database.DB.First(order, id).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	order.ProductRefer = orderBody.ProductId
	order.UserRefer = orderBody.UserId

	if err = database.Database.DB.Save(order).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("order updated")
}

func DeleteOrder(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	order := &model.Order{}
	if err = database.Database.DB.Delete(order, id).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON("order deleted successfully")
}
