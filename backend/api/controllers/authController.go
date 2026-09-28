package controllers

import (
	"context"
	"os"
	"time"

	"github.com/RomeyKung/social/database"
	"github.com/RomeyKung/social/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// Register
// @Summary Register a new user
// @Description Register a new user with the provided information
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "User information"
// @Success 201 {object} models.UserModel "Successfully registered user"
// @Failure 400 {object} map[string]interface{}
// @Router /user/signup [post]
func Register(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.CreateUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	var existingUser models.UserModel
	checkUser := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&existingUser)

	if checkUser == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user already exists",
		})
	}

	if checkUser != mongo.ErrNoDocuments {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot check existing user",
			"details": checkUser.Error(),
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot hash password",
			"details": err.Error(),
		})
	}

	newUser := models.UserModel{
		Name:      body.FirstName + " " + body.LastName,
		Email:     body.Email,
		Password:  string(hashedPassword),
		Followers: make([]string, 0),
		Following: make([]string, 0),
	}

	result, err := UserSchema.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot create user",
			"details": err.Error(),
		})
	}

	//get the new user
	var createdUser models.UserModel
	query := bson.D{{Key: "_id", Value: result.InsertedID}}

	UserSchema.FindOne(ctx, query).Decode(&createdUser)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": createdUser.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	})

	JwtSecret := os.Getenv("JWT_SECRET")

	token, err := claims.SignedString([]byte(JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot generate token",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user":  models.ToUserResponse(createdUser),
		"token": token,
	})

	// return c.Status(fiber.StatusCreated).JSON(newUser)
}

// Login
// @Summary Login a user
// @Description Login a user with the provided email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.LoginUser true "User credentials"
// @Success 200 {object} models.UserModel "Successfully logged in user"
// @Failure 400 {object} map[string]interface{}
// @Router /user/signin [post]
func Login(c *fiber.Ctx) error {

	var UserSchema = database.DB.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var body models.LoginUser
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Cannot parse JSON",
			"details": err.Error(),
		})
	}

	var user models.UserModel
	checkEmail := UserSchema.FindOne(ctx, bson.D{{Key: "email", Value: body.Email}}).Decode(&user)

	if checkEmail == mongo.ErrNoDocuments {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if checkEmail != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot check existing user",
			"details": checkEmail.Error(),
		})
	}

	checkPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if checkPassword != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expiration time
	})

	JwtSecret := os.Getenv("JWT_SECRET")

	token, err := claims.SignedString([]byte(JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Cannot generate token",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  models.ToUserResponse(user),
		"token": token,
	})

	// return c.Status(fiber.StatusCreated).JSON(newUser)
}
