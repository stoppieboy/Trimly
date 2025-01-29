package routes

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stoppieboy/trimly/database"
	"github.com/stoppieboy/trimly/helpers"
)

type Request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

type Response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(Request)

	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	r := database.CreateClient(0)
	defer r.Close()

	// implement rate limiting

	r2 := database.CreateClient(1)
	fmt.Println(r2.Ping(database.Ctx).Result())
	defer r2.Close()

	value, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second).Err()
	} else if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Some error"})
	} else {
		valInt, _ := strconv.Atoi(value)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate limit exceeded",
				"rate_limit_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	// check if the input is an actual URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Domain error"})
	}

	// enforce https, SSL
	body.URL = helpers.EnforceHTTPS(body.URL)

	var id string

	// generate a random URL if no custom short is provided
	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	// check if the custom short is already taken
	if checkShortTaken(r, id) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Custom short already taken"})
	}

	// check if the expiry is valid
	if body.Expiry == 0 {
		body.Expiry = 24
	}

	err = r.Set(database.Ctx, id, body.URL, body.Expiry*3600*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unable to connect to server"})
	}
	resp := Response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + id,
		Expiry:          body.Expiry,
		XRateRemaining:  int(r2.TTL(database.Ctx, c.IP()).Val()),
		XRateLimitReset: r2.TTL(database.Ctx, c.IP()).Val() / time.Nanosecond / time.Minute,
	}

	r2.Decr(database.Ctx, c.IP())

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func checkShortTaken(r *redis.Client,short string) bool {

	val, _ := r.Get(database.Ctx, short).Result()
	if val != "" {
		return true
	}
	return false
}
