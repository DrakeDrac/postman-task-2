package attendance

import (
	"strconv"
	"time"

	"postman-task/internal/core"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AttendanceHandler struct {
	db *gorm.DB
}

// Creates new handler
func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
	return &AttendanceHandler{
		db: db,
	}
}

func (h *AttendanceHandler) MarkAttendance(c *gin.Context) {
	var data struct {
		StudentID uint   `json:"student_id"`
		Date      string `json:"date"`
		Present   bool   `json:"present"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	// Get user ID
	markerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authorized"})
		return
	}

	// Parse date
	date, err := time.Parse("2006-01-02", data.Date)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Check if attendance exists
	var att core.Attendance
	result := h.db.Where("student_id = ? AND date = ?", data.StudentID, date).First(&att)

	if result.Error == nil {
		// Update existing
		att.Present = data.Present
		att.MarkedBy = markerID.(uint)
		h.db.Save(&att)
	} else if result.Error == gorm.ErrRecordNotFound {
		// Create new
		att = core.Attendance{
			StudentID: data.StudentID,
			Date:      date,
			Present:   data.Present,
			MarkedBy:  markerID.(uint),
		}
		h.db.Create(&att)
	} else {
		// Error
		c.JSON(500, gin.H{"error": "Database error"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Attendance marked",
		"id":      att.ID,
	})
}

// Gets attendance stats for a student
func (h *AttendanceHandler) GetAttendanceStats(c *gin.Context) {
	// Get student id from url
	studentID := c.Param("student_id")

	// Get current month's attendance
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)

	// Count present days this month
	var presentDays int64
	h.db.Model(&core.Attendance{}).
		Where("student_id = ? AND present = ? AND date >= ?", studentID, true, firstOfMonth).
		Count(&presentDays)

	// Total school days this month
	totalDays := now.Day()
	weekendDays := (now.Day() / 7) * 2
	totalSchoolDays := totalDays - weekendDays
	if totalSchoolDays < 0 {
		totalSchoolDays = 0
	}

	// Calculate percentage
	var percentage float64
	if totalSchoolDays > 0 {
		percentage = float64(presentDays) / float64(totalSchoolDays) * 100
	}

	c.JSON(200, gin.H{
		"student_id":   studentID,
		"present_days": presentDays,
		"total_days":   totalSchoolDays,
		"percentage":   percentage,
	})
}

// Gets attendance history
func (h *AttendanceHandler) GetAttendanceHistory(c *gin.Context) {
	// Get student id from url
	studentID := c.Param("student_id")

	// Get page number from query
	page := 1
	p := c.Query("page")
	if p != "" {
		pn, err := strconv.ParseInt(p, 10, 32) // parse int in base 10
		if err == nil && pn > 0 {
			page = int(pn)
		}
	}
	pageSize := 30 // Default page size
	offset := (page - 1) * pageSize

	// Get attendance records
	var records []core.Attendance
	h.db.Where("student_id = ?", studentID).
		Order("date DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records)

	// Get total count
	var total int64
	h.db.Model(&core.Attendance{}).Where("student_id = ?", studentID).Count(&total)

	c.JSON(200, gin.H{
		"data":  records,
		"page":  page,
		"total": total,
	})
}
