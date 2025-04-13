package file

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

func Upload(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		log.Print(err)

		return err
	}

	if err := c.SaveFile(file, fmt.Sprintf("./tmp/%s", file.Filename)); err != nil {
		log.Print(err)

		return err
	}

	return c.SendStatus(fiber.StatusOK)
}

func Download(c *fiber.Ctx) error {
	if err := c.Download(fmt.Sprintf("./tmp/%s", c.Params("filename")), c.Params("filename")); err != nil {
		log.Print(err)

		return err
	}

	return nil
}
