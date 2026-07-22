package jobmanagement

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, handler *Handler) {
	// Job Titles (9.1) — CRUD + Nested Subs
	titles := rg.Group("/job-management/titles")
	{
		titles.POST("", handler.CreateJobTitle)
		titles.GET("", handler.ListJobTitles)
		titles.GET("/:id", handler.GetJobTitleByID)
		titles.PUT("/:id", handler.UpdateJobTitle)
		titles.DELETE("/:id", handler.DeleteJobTitle)

		// Nested Job Title Subs (9.2)
		titles.POST("/:id/subs", handler.CreateJobTitleSub)
		titles.GET("/:id/subs", handler.ListJobTitleSubs)
		titles.GET("/:id/subs/:subId", handler.GetJobTitleSubByID)
		titles.PUT("/:id/subs/:subId", handler.UpdateJobTitleSub)
		titles.DELETE("/:id/subs/:subId", handler.DeleteJobTitleSub)
	}

	// Job Values (9.3)
	values := rg.Group("/job-management/values")
	{
		values.POST("", handler.CreateJobValue)
		values.GET("", handler.ListJobValues)
		values.GET("/:id", handler.GetJobValueByID)
		values.PUT("/:id", handler.UpdateJobValue)
		values.DELETE("/:id", handler.DeleteJobValue)
	}

	// Job Objectives (9.4)
	objectives := rg.Group("/job-management/objectives")
	{
		objectives.POST("", handler.CreateJobObjective)
		objectives.GET("", handler.ListJobObjectives)
		objectives.GET("/:id", handler.GetJobObjectiveByID)
		objectives.PUT("/:id", handler.UpdateJobObjective)
		objectives.DELETE("/:id", handler.DeleteJobObjective)
	}

	// Job Identifications (9.5)
	identifications := rg.Group("/job-management/identifications")
	{
		identifications.POST("", handler.CreateJobIdentification)
		identifications.GET("", handler.ListJobIdentifications)
		identifications.GET("/:id", handler.GetJobIdentificationByID)
		identifications.PUT("/:id", handler.UpdateJobIdentification)
		identifications.DELETE("/:id", handler.DeleteJobIdentification)
	}

	// Job Responsibilities (9.6)
	responsibilities := rg.Group("/job-management/responsibilities")
	{
		responsibilities.POST("", handler.CreateJobResponsibility)
		responsibilities.GET("", handler.ListJobResponsibilities)
		responsibilities.GET("/:id", handler.GetJobResponsibilityByID)
		responsibilities.PUT("/:id", handler.UpdateJobResponsibility)
		responsibilities.DELETE("/:id", handler.DeleteJobResponsibility)
	}

	// Job Education Experiences (9.7)
	education := rg.Group("/job-management/education-experiences")
	{
		education.POST("", handler.CreateJobEducationExperience)
		education.GET("", handler.ListJobEducationExperiences)
		education.GET("/:id", handler.GetJobEducationExperienceByID)
		education.PUT("/:id", handler.UpdateJobEducationExperience)
		education.DELETE("/:id", handler.DeleteJobEducationExperience)
	}

	// Job HR Authorities (9.8)
	hrAuth := rg.Group("/job-management/hr-authorities")
	{
		hrAuth.POST("", handler.CreateJobHRAuthority)
		hrAuth.GET("", handler.ListJobHRAuthorities)
		hrAuth.GET("/:id", handler.GetJobHRAuthorityByID)
		hrAuth.PUT("/:id", handler.UpdateJobHRAuthority)
		hrAuth.DELETE("/:id", handler.DeleteJobHRAuthority)
	}

	// Job Operational Authorities (9.9)
	opAuth := rg.Group("/job-management/operational-authorities")
	{
		opAuth.POST("", handler.CreateJobOperationalAuthority)
		opAuth.GET("", handler.ListJobOperationalAuthorities)
		opAuth.GET("/:id", handler.GetJobOperationalAuthorityByID)
		opAuth.PUT("/:id", handler.UpdateJobOperationalAuthority)
		opAuth.DELETE("/:id", handler.DeleteJobOperationalAuthority)
	}

	// Job Working Activities (9.10)
	activities := rg.Group("/job-management/working-activities")
	{
		activities.POST("", handler.CreateJobWorkingActivity)
		activities.GET("", handler.ListJobWorkingActivities)
		activities.GET("/:id", handler.GetJobWorkingActivityByID)
		activities.PUT("/:id", handler.UpdateJobWorkingActivity)
		activities.DELETE("/:id", handler.DeleteJobWorkingActivity)
	}

	// Job Working Risks (9.11)
	risks := rg.Group("/job-management/working-risks")
	{
		risks.POST("", handler.CreateJobWorkingRisk)
		risks.GET("", handler.ListJobWorkingRisks)
		risks.GET("/:id", handler.GetJobWorkingRiskByID)
		risks.PUT("/:id", handler.UpdateJobWorkingRisk)
		risks.DELETE("/:id", handler.DeleteJobWorkingRisk)
	}

	// Job Relationships (9.12)
	relationships := rg.Group("/job-management/relationships")
	{
		relationships.POST("", handler.CreateJobRelationship)
		relationships.GET("", handler.ListJobRelationships)
		relationships.GET("/:id", handler.GetJobRelationshipByID)
		relationships.PUT("/:id", handler.UpdateJobRelationship)
		relationships.DELETE("/:id", handler.DeleteJobRelationship)
	}

	// Job Subordinate Controls (9.13)
	subordinates := rg.Group("/job-management/subordinate-controls")
	{
		subordinates.POST("", handler.CreateJobSubordinateControl)
		subordinates.GET("", handler.ListJobSubordinateControls)
		subordinates.GET("/:id", handler.GetJobSubordinateControlByID)
		subordinates.PUT("/:id", handler.UpdateJobSubordinateControl)
		subordinates.DELETE("/:id", handler.DeleteJobSubordinateControl)
	}

	// Job Assets (9.14)
	assets := rg.Group("/job-management/assets")
	{
		assets.POST("", handler.CreateJobAsset)
		assets.GET("", handler.ListJobAssets)
		assets.GET("/:id", handler.GetJobAssetByID)
		assets.PUT("/:id", handler.UpdateJobAsset)
		assets.DELETE("/:id", handler.DeleteJobAsset)
	}

	// Job Financials (9.15)
	financials := rg.Group("/job-management/financials")
	{
		financials.POST("", handler.CreateJobFinancial)
		financials.GET("", handler.ListJobFinancials)
		financials.GET("/:id", handler.GetJobFinancialByID)
		financials.PUT("/:id", handler.UpdateJobFinancial)
		financials.DELETE("/:id", handler.DeleteJobFinancial)
	}

	// Job Potency Competencies (9.16)
	potency := rg.Group("/job-management/potency-competencies")
	{
		potency.POST("", handler.CreateJobPotencyCompetency)
		potency.GET("", handler.ListJobPotencyCompetencies)
		potency.GET("/:id", handler.GetJobPotencyCompetencyByID)
		potency.PUT("/:id", handler.UpdateJobPotencyCompetency)
		potency.DELETE("/:id", handler.DeleteJobPotencyCompetency)
	}

	// Job Scores (9.17) — scoped under organization
	scores := rg.Group("/job-management/scores")
	{
		scores.GET("", handler.ListJobScores)
		scores.GET("/org/:orgId", handler.GetJobScoreByOrganization)
		scores.PUT("/org/:orgId", handler.UpsertJobScore)
	}

	// Job Competency Groups (9.18)
	compGroups := rg.Group("/job-management/competency-groups")
	{
		compGroups.POST("", handler.CreateJobCompetencyGroup)
		compGroups.GET("", handler.ListJobCompetencyGroups)
		compGroups.GET("/:id", handler.GetJobCompetencyGroupByID)
		compGroups.PUT("/:id", handler.UpdateJobCompetencyGroup)
		compGroups.DELETE("/:id", handler.DeleteJobCompetencyGroup)
	}
}
