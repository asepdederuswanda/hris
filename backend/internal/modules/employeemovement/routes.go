package employeemovement

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Employee Movement ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/employee-movements
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	em := rg.Group("/employee-movements")
	{
		// Employee Movements
		em.POST("/movements", handler.CreateMovement)
		em.GET("/movements", handler.ListMovements)
		em.GET("/movements/:id", handler.GetMovementByID)
		em.PUT("/movements/:id", handler.UpdateMovement)
		em.DELETE("/movements/:id", handler.DeleteMovement)
		em.POST("/movements/:id/approve", handler.ApproveMovement)
		em.POST("/movements/:id/execute", handler.ExecuteMovement)
		em.POST("/movements/:id/cancel", handler.CancelMovement)

		// Movements by Employee
		em.GET("/employees/:employeeId/movements", handler.ListMovementsByEmployee)

		// Employee Contracts
		em.POST("/contracts", handler.CreateContract)
		em.GET("/contracts", handler.ListContracts)
		em.GET("/contracts/:id", handler.GetContractByID)
		em.PUT("/contracts/:id", handler.UpdateContract)
		em.DELETE("/contracts/:id", handler.DeleteContract)

		// Contracts by Employee
		em.GET("/employees/:employeeId/contracts", handler.ListContractsByEmployee)
	}
}
