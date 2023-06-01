package handler

import (
	"fmt"
	"gptube/services"

	"github.com/gofiber/fiber/v2"
)

func HomeHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "GPTube api",
	})
}

func ChatGPT(c *fiber.Ctx) error {
	chatInput := struct {
		Input string `json:"input"`
	}{Input: ""}

	if err := c.BodyParser(&chatInput); err != nil {
		return err
	}

	resp, err := services.Chat(chatInput.Input)
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return c.JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"output":  fmt.Sprintf("Ans: %s", resp.Choices[0].Message.Content),
		"results": resp,
	})
}
