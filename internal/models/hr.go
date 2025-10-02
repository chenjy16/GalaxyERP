package models

import (
	"time"
)

// Employee 员工模型
type Employee struct {
	BaseModel
	Code             string     `json:"code" gorm:"uniqueIndex;size:50;not null"`
	FirstName        string     `json:"first_name" gorm:"size:100;not null"`
	LastName         string     `json:"last_name" gorm:"size:100;not null"`
	FullName         string     `json:"full_name" gorm:"size:255;not null"`
	Email            string     `json:"email" gorm:"uniqueIndex;size:255;not null"`
	Phone            string     `json:"phone" gorm:"size:20"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Gender           string     `json:"gender" gorm:"size:10"`
	HireDate         *time.Time `json:"hire_date,omitempty"`
	DepartmentID     *uint      `json:"department_id,omitempty" gorm:"index"`
	PositionID       *uint      `json:"position_id,omitempty" gorm:"index"`
	ManagerID        *uint      `json:"manager_id,omitempty" gorm:"index"`
	Status           string     `json:"status" gorm:"size:50;default:'active';index"`
	EmergencyContact string     `json:"emergency_contact" gorm:"size:255"`
	IDNumber         string     `json:"id_number" gorm:"size:100"`
	Address          string     `json:"address" gorm:"type:text"`

	// 关联
	Department   *Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Position     *Position   `json:"position,omitempty" gorm:"foreignKey:PositionID"`
	Manager      *Employee   `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Subordinates []Employee  `json:"subordinates,omitempty" gorm:"foreignKey:ManagerID"`
}

// Position 职位模型
type Position struct {
	CodeModel
	DepartmentID uint `json:"department_id" gorm:"index;not null"`

	// 关联
	Department Department `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Employees  []Employee `json:"employees,omitempty" gorm:"foreignKey:PositionID"`
}

// Attendance 考勤模型
type Attendance struct {
	BaseModel
	EmployeeID   uint       `json:"employee_id" gorm:"index;not null"`
	Date         time.Time  `json:"date" gorm:"index;not null"`
	CheckInTime  *time.Time `json:"check_in_time,omitempty"`
	CheckOutTime *time.Time `json:"check_out_time,omitempty"`
	Status       string     `json:"status" gorm:"size:50;default:'present';index"` // present, absent, late, early_leave
	Notes        string     `json:"notes,omitempty" gorm:"type:text"`

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// PerformanceReview 绩效评估模型
type PerformanceReview struct {
	BaseModel
	EmployeeID   uint      `json:"employee_id" gorm:"index;not null"`
	ReviewerID   uint      `json:"reviewer_id" gorm:"index;not null"`
	ReviewDate   time.Time `json:"review_date" gorm:"index;not null"`
	ReviewPeriod string    `json:"review_period" gorm:"size:100;not null"` // Q1-2024, 2024-H1, etc.
	Score        float64   `json:"score" gorm:"not null"`
	Comments     string    `json:"comments" gorm:"type:text"`
	Status       string    `json:"status" gorm:"size:50;default:'draft';index"` // draft, submitted, approved, finalized

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Reviewer Employee `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`
}

// Skill 技能模型
type Skill struct {
	BaseModel
	Name        string `json:"name" gorm:"size:255;not null;uniqueIndex"`
	Description string `json:"description,omitempty" gorm:"type:text"`
	Category    string `json:"category" gorm:"size:100;not null;index"` // technical, soft, language, certification
	IsActive    bool   `json:"is_active" gorm:"default:true;index"`

	// 关联
	EmployeeSkills []EmployeeSkill `json:"employee_skills,omitempty" gorm:"foreignKey:SkillID"`
}

// EmployeeSkill 员工技能模型
type EmployeeSkill struct {
	BaseModel
	EmployeeID   uint       `json:"employee_id" gorm:"index;not null"`
	SkillID      uint       `json:"skill_id" gorm:"index;not null"`
	Level        string     `json:"level" gorm:"size:50;not null;index"` // beginner, intermediate, advanced, expert
	Years        float64    `json:"years" gorm:"default:0"`
	LastAssessed *time.Time `json:"last_assessed,omitempty"`

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	Skill    Skill    `json:"skill,omitempty" gorm:"foreignKey:SkillID"`
}

// Leave 请假模型
type Leave struct {
	BaseModel
	EmployeeID uint       `json:"employee_id" gorm:"index;not null"`
	LeaveType  string     `json:"leave_type" gorm:"size:50;not null"`
	StartDate  time.Time  `json:"start_date" gorm:"index;not null"`
	EndDate    time.Time  `json:"end_date" gorm:"index;not null"`
	Days       float64    `json:"days" gorm:"not null"`
	Reason     string     `json:"reason" gorm:"type:text;not null"`
	Status     string     `json:"status" gorm:"size:50;default:'pending';index"` // pending, approved, rejected, cancelled
	ApprovedBy *uint      `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// 关联
	Employee     Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	ApprovedUser *Employee `json:"approved_user,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// Payroll 工资单模型
type Payroll struct {
	BaseModel
	EmployeeID      uint       `json:"employee_id" gorm:"index;not null"`
	PayPeriodStart  time.Time  `json:"pay_period_start" gorm:"index;not null"`
	PayPeriodEnd    time.Time  `json:"pay_period_end" gorm:"index;not null"`
	BasicSalary     float64    `json:"basic_salary" gorm:"not null"`
	OvertimePay     float64    `json:"overtime_pay" gorm:"default:0"`
	Allowance       float64    `json:"allowance" gorm:"default:0"`
	Bonus           float64    `json:"bonus" gorm:"default:0"`
	Deductions      float64    `json:"deductions" gorm:"default:0"`
	SocialInsurance float64    `json:"social_insurance" gorm:"default:0"`
	HousingFund     float64    `json:"housing_fund" gorm:"default:0"`
	Tax             float64    `json:"tax" gorm:"default:0"`
	NetPay          float64    `json:"net_pay" gorm:"not null"`
	Status          string     `json:"status" gorm:"size:50;default:'draft';index"` // draft, confirmed, paid
	PaidAt          *time.Time `json:"paid_at,omitempty"`

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// Training 培训模型
type Training struct {
	BaseModel
	Code            string    `json:"code" gorm:"uniqueIndex;size:50;not null"`
	Name            string    `json:"name" gorm:"size:255;not null"`
	Description     string    `json:"description,omitempty" gorm:"type:text"`
	StartDate       time.Time `json:"start_date" gorm:"index;not null"`
	EndDate         time.Time `json:"end_date" gorm:"index;not null"`
	Trainer         string    `json:"trainer,omitempty" gorm:"size:100"`
	Location        string    `json:"location,omitempty" gorm:"size:255"`
	Cost            float64   `json:"cost,omitempty" gorm:"default:0"`
	Currency        string    `json:"currency" gorm:"size:10;default:'CNY'"`
	MaxParticipants int       `json:"max_participants,omitempty" gorm:"default:0"`
	Status          string    `json:"status" gorm:"size:50;default:'planned';index"` // planned, ongoing, completed, cancelled

	// 关联
	Participants []TrainingParticipant `json:"participants,omitempty" gorm:"foreignKey:ProgramID"`
}

// TrainingParticipant 培训参与者模型
type TrainingParticipant struct {
	BaseModel
	EmployeeID     uint      `json:"employee_id" gorm:"index;not null"`
	ProgramID      uint      `json:"program_id" gorm:"index;not null"`
	EnrollmentDate time.Time `json:"enrollment_date" gorm:"index;not null"`
	Attendance     string    `json:"attendance" gorm:"size:50;default:'enrolled';index"` // enrolled, attended, completed, cancelled
	Score          *float64  `json:"score,omitempty"`
	Certificate    string    `json:"certificate,omitempty" gorm:"size:500"`
	Status         string    `json:"status" gorm:"size:50;default:'enrolled';index"` // enrolled, attended, completed, cancelled

	// 关联
	Training Training `json:"training,omitempty" gorm:"foreignKey:ProgramID"`
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}

// OvertimeRecord 加班记录模型
type OvertimeRecord struct {
	BaseModel
	EmployeeID uint       `json:"employee_id" gorm:"index;not null"`
	Date       time.Time  `json:"date" gorm:"index;not null"`
	Hours      float64    `json:"hours" gorm:"not null"`
	Reason     string     `json:"reason" gorm:"type:text;not null"`
	Status     string     `json:"status" gorm:"size:50;default:'pending';index"` // pending, approved, rejected
	ApprovedBy *uint      `json:"approved_by,omitempty" gorm:"index"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`

	// 关联
	Employee     Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	ApprovedUser *Employee `json:"approved_user,omitempty" gorm:"foreignKey:ApprovedBy"`
}

// PerformanceGoal 绩效目标模型
type PerformanceGoal struct {
	BaseModel
	EmployeeID  uint      `json:"employee_id" gorm:"index;not null"`
	Name        string    `json:"name" gorm:"size:255;not null"`
	Description string    `json:"description" gorm:"type:text"`
	StartDate   time.Time `json:"start_date" gorm:"index;not null"`
	EndDate     time.Time `json:"end_date" gorm:"index;not null"`
	TargetValue float64   `json:"target_value" gorm:"not null"`
	ActualValue float64   `json:"actual_value" gorm:"default:0"`
	Weight      float64   `json:"weight" gorm:"default:1"`                      // 权重
	Status      string    `json:"status" gorm:"size:50;default:'active';index"` // active, completed, cancelled

	// 关联
	Employee Employee `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
}
