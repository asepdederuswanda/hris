package competency

import "github.com/gin-gonic/gin"

// RegisterRoutes mendaftarkan semua endpoint Competency ke router group tenant.
// Semua endpoint di bawah /api/v1/tenant/competency
func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	comp := rg.Group("/competency")
	{
		// Competencies (master)
		comp.POST("/competencies", handler.CreateCompetency)
		comp.GET("/competencies", handler.ListCompetencies)
		comp.GET("/competencies/:id", handler.GetCompetencyByID)
		comp.PUT("/competencies/:id", handler.UpdateCompetency)
		comp.DELETE("/competencies/:id", handler.DeleteCompetency)

		// Competence Values (legacy)
		comp.POST("/competence-values", handler.CreateCompetenceValue)
		comp.GET("/competence-values", handler.ListCompetenceValues)
		comp.GET("/competence-values/:id", handler.GetCompetenceValueByID)
		comp.PUT("/competence-values/:id", handler.UpdateCompetenceValue)
		comp.DELETE("/competence-values/:id", handler.DeleteCompetenceValue)

		// Competency Values (structured)
		comp.POST("/values", handler.CreateCompetencyValue)
		comp.GET("/values", handler.ListCompetencyValues)
		comp.GET("/values/:id", handler.GetCompetencyValueByID)
		comp.PUT("/values/:id", handler.UpdateCompetencyValue)
		comp.DELETE("/values/:id", handler.DeleteCompetencyValue)

		// Competency Events
		comp.POST("/events", handler.CreateCompetencyEvent)
		comp.GET("/events", handler.ListCompetencyEvents)
		comp.GET("/events/:id", handler.GetCompetencyEventByID)
		comp.PUT("/events/:id", handler.UpdateCompetencyEvent)
		comp.DELETE("/events/:id", handler.DeleteCompetencyEvent)

		// Competency Event Targets
		comp.POST("/event-targets", handler.CreateCompetencyEventTarget)
		comp.GET("/event-targets", handler.ListCompetencyEventTargets)
		comp.GET("/event-targets/:id", handler.GetCompetencyEventTargetByID)
		comp.PUT("/event-targets/:id", handler.UpdateCompetencyEventTarget)
		comp.DELETE("/event-targets/:id", handler.DeleteCompetencyEventTarget)

		// Competency Scores
		comp.POST("/scores", handler.CreateCompetencyScore)
		comp.GET("/scores", handler.ListCompetencyScores)
		comp.GET("/scores/:id", handler.GetCompetencyScoreByID)
		comp.PUT("/scores/:id", handler.UpdateCompetencyScore)
		comp.DELETE("/scores/:id", handler.DeleteCompetencyScore)

		// Competency Score Details
		comp.POST("/score-details", handler.CreateCompetencyScoreDetail)
		comp.GET("/scores/:scoreId/details", handler.ListCompetencyScoreDetails)
		comp.GET("/score-details/:id", handler.GetCompetencyScoreDetailByID)
		comp.PUT("/score-details/:id", handler.UpdateCompetencyScoreDetail)
		comp.DELETE("/score-details/:id", handler.DeleteCompetencyScoreDetail)
	}
}
