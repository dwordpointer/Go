package controllers

import (
	"main/database"
	"main/middleware"
	"main/models"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] == "" || data["passwordConfirm"] == "" || data["email"] == "" || data["firstName"] == "" || data["lastName"] == "" {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Boş veri girilemez.",
		})
	}

	if data["password"] != data["passwordConfirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Şifreler Eşleşmiyor.",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	isUser := models.User{
		Email: data["email"],
	}
	
	database.DB.Where("email = ?", data["email"]).First(&isUser)
	if isUser.Email == "" {
		return c.JSON(fiber.Map{
			"message": "Bu eMail zaten mevcut!",
		})
	}

	user := models.User{
		Firstname: data["firstName"],
		LastName:  data["lastName"],
		Email:     data["email"],
		Password:  password,
	}
	database.DB.Create(&user)
	return c.JSON(fiber.Map{
		"message": "Kayıt Başarılı",
	})
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.Id == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Kullanıcı bulunamadı!",
		})
	}
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Şifre Yanlış!",
		})
	}
	claims := jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("asgsagqwr2125asgasxbxz"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"jwt":     token,
		"message": "Login Success",
	})
}

func User(c *fiber.Ctx) error {
	var user models.User

	deg := c.Locals("claim")
	details := deg.(*middleware.Claims)

	database.DB.Where("id = ?", details.Issuer).First(&user)

	return c.JSON(fiber.Map{
		"id":    user.Id,
		"email": user.Email,
	})
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "Logout Success",
	})
}
