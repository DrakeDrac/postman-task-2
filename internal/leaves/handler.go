package leaves

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"postman-task/internal/core"
	email "postman-task/internal/notifications"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LeaveHandler struct {
	db *gorm.DB
}

func NewLeaveHandler(db *gorm.DB) *LeaveHandler {
	return &LeaveHandler{
		db: db,
	}
}

// Handles leave application
func (h *LeaveHandler) ApplyLeave(c *gin.Context) {
	// Get user id from gin context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authorized"})
		return
	}

	// Get leave data
	var data struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
		Reason    string `json:"reason"`
		Type      string `json:"leave_type"`
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	// Parse dates
	start, err1 := time.Parse("2006-01-02", data.StartDate)
	end, err2 := time.Parse("2006-01-02", data.EndDate)

	if err1 != nil || err2 != nil {
		c.JSON(400, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Create leave request
	leave := core.LeaveRequest{
		StudentID: userID.(uint),
		StartDate: start,
		EndDate:   end,
		Reason:    data.Reason,
		LeaveType: data.Type,
		Status:    "pending",
	}

	// Save to database
	result := h.db.Create(&leave)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "Could not save leave request"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Leave request submitted",
		"id":      leave.ID,
	})
}

// Gets all leaves for current user
func (h *LeaveHandler) GetMyLeaves(c *gin.Context) {
	// Get user id from gin context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authorized"})
		return
	}

	// Get pagination parameters
	page := 1
	p := c.Query("page")
	if p != "" {
		pn, err := strconv.ParseInt(p, 10, 32) // parse int in base 10
		if err == nil && pn > 0 {
			page = int(pn)
		}
	}
	pageSize := 10

	// Get total count
	var total int64
	h.db.Model(&core.LeaveRequest{}).Where("student_id = ?", userID).Count(&total)

	// Get page of leaves
	var leaves []core.LeaveRequest
	offset := (page - 1) * pageSize
	h.db.Where("student_id = ?", userID).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&leaves)

	// Return paginated result
	result := core.PageResult{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    leaves,
	}

	c.JSON(200, result)
}

// Handles both approval and rejection of leave requests
func (h *LeaveHandler) HandleLeaveAction(c *gin.Context) {
	// Get leave id and action from url
	leaveID := c.Param("id")
	action := c.Param("action")
	fmt.Println(action)

	// Check action
	if action != "approve" && action != "reject" {
		c.JSON(400, gin.H{"error": "Invalid action. Must be 'approve' or 'reject'"})
		return
	}

	var data struct {
		Remarks *string `json:"remarks,omitempty"`
	}

	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad request"})
		return
	}

	// Get approver id from gin context
	approverID, exists := c.Get("user_id")
	if !exists {
		c.JSON(401, gin.H{"error": "Not authorized"})
		return
	}

	// Find leave request
	var leave core.LeaveRequest
	err = h.db.First(&leave, leaveID).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Leave not found"})
		return
	}

	// Update leave status
	actionText := action + "d" // "approved" or "rejected"
	// Ensure the status is set correctly (handle both "approve" and "reject" actions)
	if action == "approve" {
		leave.Status = "approved"
	} else if action == "reject" {
		leave.Status = "rejected"
	}
	if data.Remarks != nil {
		leave.Remarks = data.Remarks
	}
	approverIDUint := approverID.(uint)
	leave.ApprovedBy = &approverIDUint

	// Save changes
	err = h.db.Save(&leave).Error
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update leave request"})
		return
	}

	// If approved, mark the student absent for every day within the leave period
	if action == "approve" {
		for d := leave.StartDate; !d.After(leave.EndDate); d = d.AddDate(0, 0, 1) {
			var count int64
			if err := h.db.Model(&core.Attendance{}).
				Where("student_id = ? AND date = ?", leave.StudentID, d).
				Count(&count).Error; err != nil {
				c.JSON(500, gin.H{"error": "Database error"})
				return
			}

			if count == 0 {
				att := core.Attendance{
					StudentID: leave.StudentID,
					Date:      d,
					Present:   false,
					MarkedBy:  approverIDUint,
				}
				if err := h.db.Create(&att).Error; err != nil {
					c.JSON(500, gin.H{"error": "Failed to create attendance"})
					return
				}
			}
		}
	}

	// Notify student via email
	var student core.User
	err = h.db.First(&student, leave.StudentID).Error
	if err == nil {
		subject := fmt.Sprintf("Leave Request #%d %s", leave.ID, actionText)

		var remarks string
		if leave.Remarks != nil {
			remarks = *leave.Remarks
		} else {
			remarks = "-"
		}

		body := "Hi " + student.Name + ",\n\n" +
			"Your leave request has been " + actionText + ".\n\n" +
			"Remarks: " + remarks + "\n\n" +
			"Regards,\nFaculty"

		// Send email in background using goroutines
		go func() {
			sendErr := email.Send(student.Email, subject, body)
			if sendErr != nil {
				log.Printf("Failed to send %s email: %v", actionText, sendErr)
			}
		}()
	}

	c.JSON(200, gin.H{
		"message": "Leave request " + actionText,
		"status":  leave.Status,
	})
}

// Gets all leave requests (only for admin/faculty/warden)
func (h *LeaveHandler) GetAllLeaves(c *gin.Context) {
	// Get pagination parameters
	page := 1
	p := c.Query("page")
	if p != "" {
		pn, err := strconv.ParseInt(p, 10, 32) // parse int in base 10
		if err == nil && pn > 0 {
			page = int(pn)
		}
	}
	pageSize := 10

	// Get total count
	var total int64
	h.db.Model(&core.LeaveRequest{}).Count(&total)

	// Get page of leaves
	var leaves []core.LeaveRequest
	offset := (page - 1) * pageSize
	h.db.Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&leaves)

	// Return paginated result
	result := core.PageResult{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Items:    leaves,
	}

	c.JSON(200, result)
}
