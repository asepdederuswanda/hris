package employee

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Employee ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/employees
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	emps := rg.Group("/employees")
	{
		// Employee CRUD
		emps.POST("", handler.Create)
		emps.GET("", handler.List)
		emps.GET("/:id", handler.GetByID)
		emps.PUT("/:id", handler.Update)
		emps.DELETE("/:id", handler.Delete)

		// Addresses
		emps.POST("/:id/addresses", handler.CreateAddress)
		emps.PUT("/:id/addresses/:addressId", handler.UpdateAddress)
		emps.DELETE("/:id/addresses/:addressId", handler.DeleteAddress)

		// Emergency Contacts
		emps.POST("/:id/emergency-contacts", handler.CreateEmergencyContact)
		emps.PUT("/:id/emergency-contacts/:contactId", handler.UpdateEmergencyContact)
		emps.DELETE("/:id/emergency-contacts/:contactId", handler.DeleteEmergencyContact)

		// Families
		emps.POST("/:id/families", handler.CreateFamily)
		emps.PUT("/:id/families/:familyId", handler.UpdateFamily)
		emps.DELETE("/:id/families/:familyId", handler.DeleteFamily)

		// Educations
		emps.POST("/:id/educations", handler.CreateEducation)
		emps.PUT("/:id/educations/:educationId", handler.UpdateEducation)
		emps.DELETE("/:id/educations/:educationId", handler.DeleteEducation)

		// Experiences
		emps.POST("/:id/experiences", handler.CreateExperience)
		emps.PUT("/:id/experiences/:experienceId", handler.UpdateExperience)
		emps.DELETE("/:id/experiences/:experienceId", handler.DeleteExperience)

		// Documents
		emps.POST("/:id/documents", handler.CreateDocument)
		emps.PUT("/:id/documents/:documentId", handler.UpdateDocument)
		emps.DELETE("/:id/documents/:documentId", handler.DeleteDocument)

		// Insurances
		emps.POST("/:id/insurances", handler.CreateInsurance)
		emps.PUT("/:id/insurances/:insuranceId", handler.UpdateInsurance)
		emps.DELETE("/:id/insurances/:insuranceId", handler.DeleteInsurance)

		// Employments
		emps.POST("/:id/employments", handler.CreateEmployment)
		emps.PUT("/:id/employments/:employmentId", handler.UpdateEmployment)
		emps.DELETE("/:id/employments/:employmentId", handler.DeleteEmployment)
	}
}
