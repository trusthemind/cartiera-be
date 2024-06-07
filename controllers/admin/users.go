package admin_controllers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"golang.org/x/crypto/bcrypt"

	"github.com/trusthemind/go-cars-app/helpers"
	"github.com/trusthemind/go-cars-app/initializers"
	"github.com/trusthemind/go-cars-app/models"
)

//	@Tags			Administration
//	@Summary		Users Administration
//	@Description	Get all users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.User
//	@Failure		400	{object}	models.Error
//	@Router			/admin/users [get]
func GetAllUsers(c *gin.Context) {
	var users []models.User
	result := initializers.DB.Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

//	@Tags			Administration
//	@Summary		Users Administration
//	@Description	Delete User by ID
//	@Accept			json
//	@Params			user_id path string "User ID"
//	@Success		200	{object}	models.Message
//	@Failure		400	{object}	models.Error
//	@Router			/admin/users/delete/:id [delete]
func DeleteUserbyID(c *gin.Context) {
	var userID = c.Param("id")
	var user models.User

	result := initializers.DB.Where("ID = ?", userID).Delete(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failde to delete user"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "User with ID " + userID + " was successfully deleted"})
}

//	@Tags			Administration
//	@Summary		Users Administration
//	@Description	Update User by ID
//	@Accept			json
//	@Params			user_id path string "User ID"
//	@Success		200	{object}	models.Message
//	@Failure		400	{object}	models.Error
//	@Router			/admin/users/update/:id [put]
func UpdateUserByID(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	userID := c.Param("id")
	user := models.User{}

	dbFields := map[string]interface{}{}
	for key, value := range requestBody {
		dbFields[key] = value
	}

	result := initializers.DB.Model(&user).Where("owner_id = ?", userID).Updates(dbFields)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "This user is not found"})
		return
	} else if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User has been successfully updated"})
}

//	@Tags		Administration
//	@Summary	Create a new user with admin account
//	@Description
//	@Accept		json
//	@Produce	json
//	@Param		request	body		models.AdminRequestRegistration	true	"User Data"
//	@Success	200		{object}	models.Message
//	@Failure	400		{object}	models.Error
//	@Router		/admin/new-user [post]
func CreateNewUser(c *gin.Context) {
	var requestBody models.AdminRequestRegistration

	if c.Bind(&requestBody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var customerID string
	if requestBody.IsCustomer == true {
		stripe.Key = os.Getenv("STRIPE_KEY")

		params := &stripe.CustomerParams{
			Email: &requestBody.Email,
			Name:  &requestBody.Name,
		}

		customer, err := customer.New(params)

		customerID = customer.ID

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create customer"})
			return
		}
		log.Print(customer)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}

	user := models.User{
		CustomerID: helpers.IsCustomerCheck(requestBody.IsCustomer, customerID),
		Name:       requestBody.Name,
		Email:      requestBody.Email,
		Password:   string(hash),
	}

	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create new user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
