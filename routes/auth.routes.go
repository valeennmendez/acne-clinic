package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	//"github.com/golang-jwt/jwt"
	"github.com/valeennmendez/api-go/connection"
	"github.com/valeennmendez/api-go/models"
	"golang.org/x/crypto/bcrypt"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func RegisterUser(c *gin.Context) {
	var userAdmin models.Admin

	if err := c.ShouldBindJSON(&userAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format: " + err.Error(),
		})
		return
	}

	if existEmail(userAdmin.Email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "The email already exists, please use another",
		})
		return
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(userAdmin.Password), bcrypt.DefaultCost)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password not hashed" + err.Error(),
		})
		return
	}

	userAdmin.Password = string(passwordHashed)
	userAdmin.Status = "pending"
	userAdmin.Role = "common"

	if err := connection.DB.Create(&userAdmin).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create UserAdmin: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "UserAdmin created succesfully, pending approval",
	})
}

func existEmail(email string) bool {
	var user models.Admin

	// Busca un registro donde el campo email sea igual al proporcionado
	if err := connection.DB.Where("email = ?", email).First(&user).Error; err != nil {
		// Si no encuentra el registro, el error es record not found
		if err.Error() == "record not found" {
			return false
		}
		// Si hay otro tipo de error, lo manejamos aquí (opcionalmente se puede registrar el error)
		return false
	}
	// Si encuentra el registro, retorna true
	return true
}

func Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json: "email"`
		Password string `json: "password"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format",
		})
		return
	}

	var user models.Admin

	if err := connection.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if user.Status != "approved" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not approved"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	session, _ := store.Get(c.Request, "session-name")
	session.Values["authenticated"] = true
	session.Values["userName"] = user.FullName
	session.Values["role"] = user.Role

	session.Options = &sessions.Options{
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}

	err := session.Save(c.Request, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login succesful",
	})

}

func GetUserName(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")

	userName, ok := session.Values["userName"].(string)
	if !ok || userName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userName": userName})
}

func GetRoleAdmin(c *gin.Context) {

	session, _ := store.Get(c.Request, "session-name")

	role, ok := session.Values["role"].(string)

	fmt.Println("rol:", role)
	fmt.Println("accediendo...")

	if !ok || role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "you do not have permissions to access this section",
		})
		return
	}

	c.JSON(http.StatusOK, role)

}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session-name")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "No se pudo establecer la session",
			})
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()

	}
}

func ValidateSession(c *gin.Context) {
	session, err := store.Get(c.Request, "session-name")
	if err != nil || session.Values["authenticated"] != true {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Session valid",
	})
}

func CloseSesion(c *gin.Context) {
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Se produjo un error al obtener la sesión",
		})
		return
	}

	// Establecer el valor de autenticación en falso y eliminar la cookie
	session.Values["authenticated"] = false
	session.Options.MaxAge = -1 // Esto asegura que la cookie de sesión sea eliminada
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteNoneMode

	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Se produjo un error al guardar la sesión",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sesión cerrada",
	})
}

func ApproveUser(c *gin.Context) {
	var user models.Admin

	userID := c.Param("id")

	fmt.Println(userID)

	if err := connection.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	user.Status = "approved"

	if err := connection.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to approve user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user approved",
	})

}

func GetAllAdmins(c *gin.Context) {
	var admins []models.Admin

	connection.DB.Find(&admins)

	c.JSON(http.StatusOK, admins)

}

func DeclineUser(c *gin.Context) {
	var user models.Admin

	userID := c.Param("id")

	fmt.Println(userID)

	if err := connection.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user not found",
		})
		return
	}

	user.Status = "decline"

	if err := connection.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to decline user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user decline",
	})

}

func RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session-name")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "no se pudo establecer la sesion",
			})
			c.Abort()
			return
		}

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		userRole, ok := session.Values["role"].(string)
		if !ok || userRole != role {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
