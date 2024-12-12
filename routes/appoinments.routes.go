package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/valeennmendez/api-go/connection"

	"github.com/valeennmendez/api-go/models"
	"gorm.io/datatypes"
)

var aviableHours = []string{
	"08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30", "12:00", "12:30",
	"13:00", "13:30", "14:00", "14:30", "15:00", "15:30", "16:00", "16:30", "17:00", "17:30",
}

type CreateAppointmentInput struct {
	PacienteID     int    `json:"pacienteid"`
	Fecha          string `json:"fecha"`
	Hora           string `json:"hora"`
	MotivoConsulta string `json:"motivoconsulta"`
	Status         string `json: "status"`
}

func GetAviableHours(c *gin.Context) {
	dataStr := c.Query("fecha")

	fmt.Println(dataStr)

	if dataStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "date query parameter is required",
		})
		return
	}

	date, err := time.Parse("2006-01-02", dataStr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format",
		})
		return
	}

	var bookedAppointments []models.Appoinment
	connection.DB.Where("fecha = ?", date).Find(&bookedAppointments)

	bookedHours := make(map[string]bool)

	for _, appointment := range bookedAppointments {
		bookedHours[appointment.Hora] = true
	}

	aviableHoursList := []string{}

	for _, hour := range aviableHours {
		if !bookedHours[hour] {
			aviableHoursList = append(aviableHoursList, hour)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"available_hours": aviableHoursList,
	})

}

func CreateAppoinment(c *gin.Context) {

	var appoinment CreateAppointmentInput

	if err := c.ShouldBindJSON(&appoinment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON format: " + err.Error(),
		})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", appoinment.Fecha)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fecha inválida, formato debe ser YYYY-MM-DD"})
		return
	}

	newAppointment := models.Appoinment{
		PacienteID:     appoinment.PacienteID,
		Fecha:          datatypes.Date(parsedDate),
		Hora:           appoinment.Hora,
		MotivoConsulta: appoinment.MotivoConsulta,
		Status:         "pending",
	}

	validHour := false

	for _, i := range aviableHours {
		if appoinment.Hora == i {
			validHour = true
			break
		}
	}

	if !validHour {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The selected time slot is not available",
		})
		return
	}

	var existingAppointment models.Appoinment

	err = connection.DB.Where("fecha = ? AND hora = ?", newAppointment.Fecha, newAppointment.Hora).First(&existingAppointment).Error

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The selected time slot is already booked",
		})
		return
	}

	var pacienteEncontrado models.Patients

	if err := connection.DB.Where("id = ?", newAppointment.PacienteID).First(&pacienteEncontrado).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Patient not found." + err.Error(),
		})
	}

	if err := connection.DB.Create(&newAppointment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create appoinment" + err.Error(),
		})
		return
	}

	/* 	appointmentData := email.AppointmentData{
	   		Name:          pacienteEncontrado.FullName,
	   		Date:          appoinment.Fecha,
	   		Time:          appoinment.Hora,
	   		Motivo:        appoinment.MotivoConsulta,
	   		Status:        appoinment.Status,
	   		AppointmentID: int(newAppointment.ID),
	   	}

	   	go func() {
	   		err = email.SendEmail(pacienteEncontrado.Email, "Appointment Confirmation", appointmentData)
	   		if err != nil {
	   			log.Printf("Failed to send confirmation email: %v", err)
	   		}
	   	}() */

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Appoinment created succesfully",
	})

}

func GetAllAppointments(c *gin.Context) {
	var appointments []models.Appoinment

	err := connection.DB.Preload("Paciente").Find(&appointments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to loading appointments data.",
		})
	}

	c.JSON(http.StatusOK, appointments)
}

func GetAppointmentsFilter(c *gin.Context) {
	var appointments []models.Appoinment

	status := c.Query("status")

	if err := connection.DB.Preload("Paciente").Where("status = ?", status).Find(&appointments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to loading appointments",
		})
		return
	}

	c.JSON(http.StatusOK, appointments)

}

func CancelAppointment(c *gin.Context) {
	var appointment models.Appoinment

	id := c.Param("id")

	if err := connection.DB.First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Appointment not found" + err.Error(),
		})
		return
	}

	/* 	if err := connection.DB.Delete(&appointment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to cancel appointment" + err.Error(),
		})
		return
	} */

	fmt.Println(appointment.Status)
	appointment.Status = "canceled"

	if err := connection.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant save to new data",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment cancelated successfully",
	})
}

func DoneAppointment(c *gin.Context) {
	var appointment models.Appoinment

	id := c.Param("id")

	if err := connection.DB.First(&appointment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Appointment not found" + err.Error(),
		})
		return
	}

	appointment.Status = "done"

	if err := connection.DB.Save(&appointment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cant save to new data",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "appointment doned successfully",
	})
}

func GetAppointmentID(c *gin.Context) {
	var appoinment models.Appoinment

	id := c.Param("id")

	if err := connection.DB.Preload("Paciente").First(&appoinment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Appointment not found",
		})
	}

	c.JSON(http.StatusAccepted, &appoinment)
}

func EditAppointment(c *gin.Context) {
	var appoinment models.Appoinment

	id := c.Param("id")

	if err := connection.DB.First(&appoinment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Patient not found",
		})
		return
	}

	if err := c.ShouldBindJSON(&appoinment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid JSON formart" + err.Error(),
		})
		return
	}

	if err := connection.DB.Save(&appoinment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to actualice patient.",
		})
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "Appointment edited succesfully",
	})

}

func NextAppointments(c *gin.Context) {

	var appointments []models.Appoinment

	timeParse := time.Now().Format("2006-01-02")

	if err := connection.DB.Where("fecha >= ? AND status = ? ", timeParse, "pending").Order("fecha ASC, hora ASC").Limit(10).Preload("Paciente").Find(&appointments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error retrieving appointments",
		})
		return
	}

	c.JSON(http.StatusOK, appointments)

}

func SearchAppointment(c *gin.Context) {
	search := c.Query("p")

	var appointment []models.Appoinment

	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "search term is required",
		})
		return
	}

	if err := connection.DB.Joins("JOIN patients ON patients.id = appoinments.paciente_id").
		Where("patients.full_name LIKE ? OR patients.dni LIKE ?", "%"+search+"%", "%"+search+"%").Preload("Paciente").
		Find(&appointment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "error searching the patient" + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, appointment)

}
