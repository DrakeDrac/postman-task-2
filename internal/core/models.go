package core

import (
	"time"

	"gorm.io/gorm"
)

// Represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Password  string         `json:"-" gorm:"not null"` // "-" means hide from json
	Role      string         `json:"role" gorm:"not null"`
	Dept      string         `json:"dept" gorm:"not null"`
	Leaves    []LeaveRequest `json:"leaves,omitempty" gorm:"foreignKey:StudentID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Represents a leave application
type LeaveRequest struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	StudentID  uint           `json:"student_id" gorm:"not null;index"`
	Student    User           `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	LeaveType  string         `json:"leave_type" gorm:"not null;check:leave_type IN ('Medical','Personal','Academic','Emergency')"`
	Reason     string         `json:"reason" gorm:"not null"`
	StartDate  time.Time      `json:"start_date" gorm:"not null"`
	EndDate    time.Time      `json:"end_date" gorm:"not null"`
	Status     string         `json:"status" gorm:"not null;default:'pending';check:status IN ('pending','approved','rejected')"`
	ApprovedBy *uint          `json:"approved_by,omitempty" gorm:"index"`
	Approver   *User          `json:"approver,omitempty" gorm:"foreignKey:ApprovedBy"`
	Remarks    *string        `json:"remarks,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}

// Represents daily attendance records
type Attendance struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	StudentID uint           `json:"student_id" gorm:"not null;index"`
	Student   User           `json:"student,omitempty" gorm:"foreignKey:StudentID"`
	Date      time.Time      `json:"date" gorm:"not null;index"`
	Present   bool           `json:"present" gorm:"not null;default:false"`
	MarkedBy  uint           `json:"marked_by" gorm:"not null"`
	Marker    User           `json:"marker,omitempty" gorm:"foreignKey:MarkedBy"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Represents attendance statistics for a student
type AttendanceStats struct {
	StudentID            uint    `json:"student_id"`
	PresentDays          int     `json:"present_days"`
	TotalDays            int     `json:"total_days"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

// Holds pagination details
type PageInfo struct {
	Page     int `json:"page"`      // Current page number
	PageSize int `json:"page_size"` // Number of items per page
}

// Represents paginated results
type PageResult struct {
	Page     int         `json:"page"`      // Current page
	PageSize int         `json:"page_size"` // Items per page
	Total    int64       `json:"total"`     // Total items
	Items    interface{} `json:"items"`     // The actual data
}

// Represents an api response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Login request body
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// Registration request body
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"`
	Dept     string `json:"dept" binding:"required,min=2"`
}

// Leave application request body
type LeaveApplicationRequest struct {
	StudentID uint   `json:"student_id" binding:"required"`
	LeaveType string `json:"leave_type" binding:"required"`
	Reason    string `json:"reason" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
}

// Leave approval request body
type LeaveApprovalRequest struct {
	Status  string  `json:"status" binding:"required"`
	Remarks *string `json:"remarks,omitempty"`
}

// Attendance marking request body
type AttendanceMarkRequest struct {
	StudentID uint   `json:"student_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	Present   bool   `json:"present"`
}
