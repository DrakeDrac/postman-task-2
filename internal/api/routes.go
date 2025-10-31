package api

import (
	"postman-task/internal/attendance"
	"postman-task/internal/auth"
	"postman-task/internal/leaves"
	"postman-task/internal/users"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup the API routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, jwt *auth.JWTManager) {
	// Create handlers
	userH := users.NewUserHandler(db, jwt)
	leaveH := leaves.NewLeaveHandler(db)
	attendanceH := attendance.NewAttendanceHandler(db)
	analyticsH := NewAnalyticsHandler(db)

	// User routes
	r.POST("/api/v1/auth/register", userH.Register)
	r.POST("/api/v1/auth/login", userH.Login)

	// Needs token
	authorized := r.Group("/api/v1")
	authorized.Use(jwt.AuthMiddleware())
	{
		// User routes
		authorized.GET("/users", jwt.AdminOnly(), userH.GetUsers)
		authorized.GET("/users/:id", userH.GetUserByID)

		// Leave routes
		authorized.POST("/leaves/apply", leaveH.ApplyLeave)
		authorized.GET("/leaves/my", leaveH.GetMyLeaves)
		authorized.GET("/leaves", jwt.FacultyOrWarden(), leaveH.GetAllLeaves)

		// Handle both approve and reject
		authorized.PUT("/leaves/:id/:action", jwt.FacultyOrWarden(), leaveH.HandleLeaveAction)

		// Attendance routes
		authorized.POST("/attendance/mark", attendanceH.MarkAttendance)
		authorized.GET("/attendance/stats/:student_id", attendanceH.GetAttendanceStats)
		authorized.GET("/attendance/history/:student_id", attendanceH.GetAttendanceHistory)

		// Admin only routes
		admin := authorized.Group("")
		admin.Use(jwt.AdminOnly())
		{
			admin.GET("/analytics/summary", analyticsH.GetSummary)
		}
	}
}
