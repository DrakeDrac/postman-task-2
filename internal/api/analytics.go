package api

import (
	"postman-task/internal/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handles basic stats
type AnalyticsHandler struct {
	db *gorm.DB
}

// Creates new handler
func NewAnalyticsHandler(db *gorm.DB) *AnalyticsHandler {
	return &AnalyticsHandler{db: db}
}

// Gets basic stats
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	// Count users by role
	var studentCount, facultyCount, wardenCount, adminCount int64
	h.db.Model(&core.User{}).Where("role = ?", "student").Count(&studentCount)
	h.db.Model(&core.User{}).Where("role = ?", "faculty").Count(&facultyCount)
	h.db.Model(&core.User{}).Where("role = ?", "warden").Count(&wardenCount)
	h.db.Model(&core.User{}).Where("role = ?", "admin").Count(&adminCount)

	// Count leave requests by status
	var pending, approved, rejected int64
	h.db.Model(&core.LeaveRequest{}).Where("status = ?", "pending").Count(&pending)
	h.db.Model(&core.LeaveRequest{}).Where("status = ?", "approved").Count(&approved)
	h.db.Model(&core.LeaveRequest{}).Where("status = ?", "rejected").Count(&rejected)

	// Get recent leaves
	var recentLeaves []core.LeaveRequest
	h.db.Preload("Student").
		Order("created_at DESC").
		Limit(10).
		Find(&recentLeaves)

	// Basic attendance stats
	var present, absent int64
	h.db.Model(&core.Attendance{}).Where("present = ?", true).Count(&present)
	h.db.Model(&core.Attendance{}).Where("present = ?", false).Count(&absent)

	c.JSON(200, gin.H{
		"users": gin.H{
			"students": studentCount,
			"faculty":  facultyCount,
			"wardens":  wardenCount,
			"admins":   adminCount,
		},
		"leaves": gin.H{
			"pending":  pending,
			"approved": approved,
			"rejected": rejected,
		},
		"attendance": gin.H{
			"present": present,
			"absent":  absent,
		},
		"recent_leaves": recentLeaves,
	})
}
