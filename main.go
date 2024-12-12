package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/valeennmendez/api-go/connection"
	"github.com/valeennmendez/api-go/models"
	"github.com/valeennmendez/api-go/routes"
)

func main() {
	connection.ConnectionDB()

	connection.DB.AutoMigrate(&models.Patients{})
	connection.DB.AutoMigrate(&models.Admin{})
	connection.DB.AutoMigrate(&models.Appoinment{})

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Servir archivos estáticos
	r.Static("/static", "./static")
	r.StaticFile("/login.html", "./pages/login.html")

	// Rutas de autenticación
	r.POST("/register", routes.RegisterUser)
	r.POST("/login", routes.Login)
	r.GET("/validate", routes.ValidateSession)
	r.POST("/logout", routes.CloseSesion)

	// Rutas protegidas
	protected := r.Group("/")
	protected.Use(routes.AuthMiddleware())
	{
		protected.StaticFile("/index.html", "./pages/index.html")
		protected.GET("/patients", routes.GetAllPatients)
		protected.GET("/patients/:id", routes.GetPatientByID)
		protected.POST("/create", routes.CreatePatient)
		protected.PUT("/edit/:id", routes.EditPatient)
		protected.DELETE("/patients/:id", routes.DeletePacients)
		protected.GET("/total-patients", routes.TotalPatientsData)
		protected.POST("/create-appointment", routes.CreateAppoinment)
		protected.GET("/appointment-today", routes.AppointmentToday)
		protected.GET("/available-hours", routes.GetAviableHours)
		protected.GET("/appointments-week", routes.AppointmentWeek)

		protected.GET("/get-username", routes.GetUserName)
		protected.GET("/next-appointments", routes.NextAppointments)
		protected.GET("/search-appointments", routes.SearchAppointment)
		protected.GET("/appointments", routes.GetAllAppointments)
		protected.GET("/appointments-filter", routes.GetAppointmentsFilter)
		protected.POST("/appointments-edit/:id", routes.EditAppointment)
		protected.GET("/appointments/:id", routes.GetAppointmentID)
		protected.GET("/search-patient", routes.SearchPatient)
		protected.DELETE("/cancel-appointment/:id", routes.CancelAppointment) // <--- DEBE ESTAR PUBLICA SI O SI.
		protected.DELETE("/done-appointment/:id", routes.DoneAppointment)     // <--- DEBE ESTAR PUBLICA SI O SI.
		protected.GET("admin-role", routes.GetRoleAdmin)
	}

	adminProtected := protected.Group("/")
	adminProtected.Use(routes.RoleMiddleware("root"))
	{
		adminProtected.POST("/approve-user/:id", routes.ApproveUser) // 2
		adminProtected.POST("/decline-user/:id", routes.DeclineUser) // 3
		adminProtected.GET("/admins", routes.GetAllAdmins)           // 1
	}

	// Ruta raíz
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusAccepted, gin.H{
			"message": "Corriendo",
		})
	})

	r.Run(":8080")
}
