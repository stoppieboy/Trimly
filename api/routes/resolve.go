package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stoppieboy/trimly/database"
)

func ResolveURL(c *fiber.Ctx) error{
	url := c.Params("url")
	r := database.CreateClient(0)
	defer r.Close()

	// check if the short mapping exists
	value, err := r.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	}else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot connect to the database"})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	// redirect to the mapped URL
	return c.Redirect(value, 301)

}