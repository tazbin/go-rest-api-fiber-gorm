package routes

import (
	"fmt"
	"rest-api/database"
	"rest-api/model"

	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func createProductResponse(product *model.Product) *Product {
	return &Product{
		ID:           product.ID,
		Name:         product.Name,
		SerialNumber: product.SerialNumber,
	}
}

func CreateProduct(c *fiber.Ctx) error {
	var product *model.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := database.Database.DB.Create(&product).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	r := createProductResponse(product)
	return c.Status(400).JSON(r)
}

func GetProducts(c *fiber.Ctx) error {
	products := []*model.Product{}

	if err := database.Database.DB.Find(&products).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	rProducts := []*Product{}
	for _, product := range products {
		rProducts = append(rProducts, createProductResponse(product))
	}
	return c.Status(200).JSON(rProducts)
}

func GetProductById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	product := &model.Product{}
	err = database.Database.DB.Find(product, id).Error

	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(500).JSON(product)

}

func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(500).JSON(err.Error())
	}

	product := &model.Product{}
	database.Database.DB.Find(product, id)

	type productUpdate struct {
		Name         string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	uProduct := &productUpdate{}
	if err := c.BodyParser(uProduct); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	product.Name = uProduct.Name
	product.SerialNumber = uProduct.SerialNumber

	if err = database.Database.DB.Save(product).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(createProductResponse(product))
}

func DeleteProductById(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	product := &model.Product{}
	if err = database.Database.DB.Delete(product, id).Error; err != nil {
		return c.Status(500).JSON(err.Error())
	}

	return c.Status(200).JSON(fmt.Sprintf("product deleted of id %v", id))
}
