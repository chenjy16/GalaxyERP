-- ============================================================================
-- GalaxyERP 人力资源模块 - PostgreSQL 建表脚本
-- 生成时间: 2025-01-01
-- 说明: 基于 /internal/models/hr.go 的结构，包含员工、职位、薪资管理
-- ============================================================================

-- ============================================================================
-- 职位管理
-- ============================================================================

-- 职位表
CREATE TABLE IF NOT EXISTS positions (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  department_id BIGINT NOT NULL,
  level INTEGER DEFAULT 1,
  min_salary DECIMAL(15,2) NULL,
  max_salary DECIMAL(15,2) NULL,
  requirements TEXT NULL,
  responsibilities TEXT NULL,
  CONSTRAINT uq_positions_code UNIQUE (code),
  CONSTRAINT fk_positions_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_positions_deleted_at ON positions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_positions_is_active ON positions (is_active);
CREATE INDEX IF NOT EXISTS idx_positions_department_id ON positions (department_id);
CREATE INDEX IF NOT EXISTS idx_positions_code ON positions (code);
CREATE INDEX IF NOT EXISTS idx_positions_level ON positions (level);

-- ============================================================================
-- 员工管理
-- ============================================================================

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  user_id BIGINT NULL,
  employee_number VARCHAR(50) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  gender VARCHAR(10) NULL,
  birth_date DATE NULL,
  id_card VARCHAR(50) NULL,
  phone VARCHAR(50) NULL,
  email VARCHAR(255) NULL,
  address TEXT NULL,
  emergency_contact VARCHAR(255) NULL,
  emergency_phone VARCHAR(50) NULL,
  hire_date DATE NOT NULL,
  probation_end_date DATE NULL,
  termination_date DATE NULL,
  department_id BIGINT NOT NULL,
  position_id BIGINT NOT NULL,
  manager_id BIGINT NULL,
  employment_type VARCHAR(20) DEFAULT 'FULL_TIME',
  status VARCHAR(20) DEFAULT 'ACTIVE',
  CONSTRAINT uq_employees_code UNIQUE (code),
  CONSTRAINT uq_employees_employee_number UNIQUE (employee_number),
  CONSTRAINT uq_employees_user_id UNIQUE (user_id),
  CONSTRAINT fk_employees_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
  CONSTRAINT fk_employees_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IF NOT EXISTS idx_employees_is_active ON employees (is_active);
CREATE INDEX IF NOT EXISTS idx_employees_user_id ON employees (user_id);
CREATE INDEX IF NOT EXISTS idx_employees_employee_number ON employees (employee_number);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees (department_id);
CREATE INDEX IF NOT EXISTS idx_employees_position_id ON employees (position_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees (manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees (status);
CREATE INDEX IF NOT EXISTS idx_employees_employment_type ON employees (employment_type);
CREATE INDEX IF NOT EXISTS idx_employees_hire_date ON employees (hire_date);

-- ============================================================================
-- 薪资管理
-- ============================================================================

-- 薪资等级表
CREATE TABLE IF NOT EXISTS salary_grades (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  min_salary DECIMAL(15,2) NOT NULL,
  max_salary DECIMAL(15,2) NOT NULL,
  CONSTRAINT uq_salary_grades_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_salary_grades_deleted_at ON salary_grades (deleted_at);
CREATE INDEX IF NOT EXISTS idx_salary_grades_is_active ON salary_grades (is_active);
CREATE INDEX IF NOT EXISTS idx_salary_grades_code ON salary_grades (code);

-- 薪资记录表
CREATE TABLE IF NOT EXISTS salary_records (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  employee_id BIGINT NOT NULL,
  salary_grade_id BIGINT NULL,
  base_salary DECIMAL(15,2) NOT NULL,
  allowances DECIMAL(15,2) DEFAULT 0,
  deductions DECIMAL(15,2) DEFAULT 0,
  overtime_pay DECIMAL(15,2) DEFAULT 0,
  bonus DECIMAL(15,2) DEFAULT 0,
  total_salary DECIMAL(15,2) NOT NULL,
  effective_date DATE NOT NULL,
  end_date DATE NULL,
  notes TEXT NULL,
  CONSTRAINT fk_salary_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_salary_records_salary_grade FOREIGN KEY (salary_grade_id) REFERENCES salary_grades(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_salary_records_deleted_at ON salary_records (deleted_at);
CREATE INDEX IF NOT EXISTS idx_salary_records_employee_id ON salary_records (employee_id);
CREATE INDEX IF NOT EXISTS idx_salary_records_salary_grade_id ON salary_records (salary_grade_id);
CREATE INDEX IF NOT EXISTS idx_salary_records_effective_date ON salary_records (effective_date);
CREATE INDEX IF NOT EXISTS idx_salary_records_end_date ON salary_records (end_date);

-- ============================================================================
-- 考勤管理
-- ============================================================================

-- 考勤记录表
CREATE TABLE IF NOT EXISTS attendances (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  employee_id BIGINT NOT NULL,
  date DATE NOT NULL,
  check_in_time TIMESTAMP NULL,
  check_out_time TIMESTAMP NULL,
  break_start_time TIMESTAMP NULL,
  break_end_time TIMESTAMP NULL,
  work_hours DECIMAL(5,2) DEFAULT 0,
  overtime_hours DECIMAL(5,2) DEFAULT 0,
  late_minutes INTEGER DEFAULT 0,
  early_leave_minutes INTEGER DEFAULT 0,
  status VARCHAR(20) DEFAULT 'NORMAL',
  notes TEXT NULL,
  CONSTRAINT fk_attendances_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT uq_attendances_employee_date UNIQUE (employee_id, date)
);
CREATE INDEX IF NOT EXISTS idx_attendances_deleted_at ON attendances (deleted_at);
CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances (employee_id);
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances (date);
CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances (status);

-- ============================================================================
-- 请假管理
-- ============================================================================

-- 请假类型表
CREATE TABLE IF NOT EXISTS leave_types (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  max_days_per_year INTEGER NULL,
  is_paid BOOLEAN DEFAULT TRUE,
  requires_approval BOOLEAN DEFAULT TRUE,
  CONSTRAINT uq_leave_types_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_leave_types_deleted_at ON leave_types (deleted_at);
CREATE INDEX IF NOT EXISTS idx_leave_types_is_active ON leave_types (is_active);
CREATE INDEX IF NOT EXISTS idx_leave_types_code ON leave_types (code);

-- 请假申请表
CREATE TABLE IF NOT EXISTS leave_requests (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  employee_id BIGINT NOT NULL,
  leave_type_id BIGINT NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  days DECIMAL(5,2) NOT NULL,
  reason TEXT NULL,
  status VARCHAR(20) DEFAULT 'PENDING',
  approved_by BIGINT NULL,
  approved_at TIMESTAMP NULL,
  rejection_reason TEXT NULL,
  CONSTRAINT fk_leave_requests_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_leave_type FOREIGN KEY (leave_type_id) REFERENCES leave_types(id) ON DELETE CASCADE,
  CONSTRAINT fk_leave_requests_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_leave_requests_deleted_at ON leave_requests (deleted_at);
CREATE INDEX IF NOT EXISTS idx_leave_requests_employee_id ON leave_requests (employee_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_leave_type_id ON leave_requests (leave_type_id);
CREATE INDEX IF NOT EXISTS idx_leave_requests_status ON leave_requests (status);
CREATE INDEX IF NOT EXISTS idx_leave_requests_start_date ON leave_requests (start_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_end_date ON leave_requests (end_date);
CREATE INDEX IF NOT EXISTS idx_leave_requests_approved_by ON leave_requests (approved_by);

-- ============================================================================
-- 绩效管理
-- ============================================================================

-- 绩效评估表
CREATE TABLE IF NOT EXISTS performance_reviews (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  employee_id BIGINT NOT NULL,
  reviewer_id BIGINT NOT NULL,
  review_period_start DATE NOT NULL,
  review_period_end DATE NOT NULL,
  overall_score DECIMAL(5,2) NULL,
  goals_achievement DECIMAL(5,2) NULL,
  competency_score DECIMAL(5,2) NULL,
  behavior_score DECIMAL(5,2) NULL,
  strengths TEXT NULL,
  areas_for_improvement TEXT NULL,
  development_plan TEXT NULL,
  status VARCHAR(20) DEFAULT 'DRAFT',
  submitted_at TIMESTAMP NULL,
  approved_at TIMESTAMP NULL,
  CONSTRAINT fk_performance_reviews_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_reviewer FOREIGN KEY (reviewer_id) REFERENCES employees(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_deleted_at ON performance_reviews (deleted_at);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_employee_id ON performance_reviews (employee_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_reviewer_id ON performance_reviews (reviewer_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_status ON performance_reviews (status);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_period ON performance_reviews (review_period_start, review_period_end);

-- ============================================================================
-- 培训管理
-- ============================================================================

-- 培训课程表
CREATE TABLE IF NOT EXISTS training_courses (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  instructor VARCHAR(255) NULL,
  duration_hours INTEGER NOT NULL,
  max_participants INTEGER NULL,
  cost DECIMAL(15,2) DEFAULT 0,
  location VARCHAR(255) NULL,
  CONSTRAINT uq_training_courses_code UNIQUE (code)
);
CREATE INDEX IF NOT EXISTS idx_training_courses_deleted_at ON training_courses (deleted_at);
CREATE INDEX IF NOT EXISTS idx_training_courses_is_active ON training_courses (is_active);
CREATE INDEX IF NOT EXISTS idx_training_courses_code ON training_courses (code);

-- 培训记录表
CREATE TABLE IF NOT EXISTS training_records (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NULL,
  updated_at TIMESTAMP NULL,
  deleted_at TIMESTAMP NULL,
  created_by BIGINT NULL,
  updated_by BIGINT NULL,
  employee_id BIGINT NOT NULL,
  course_id BIGINT NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NULL,
  status VARCHAR(20) DEFAULT 'ENROLLED',
  score DECIMAL(5,2) NULL,
  certificate_issued BOOLEAN DEFAULT FALSE,
  certificate_number VARCHAR(100) NULL,
  notes TEXT NULL,
  CONSTRAINT fk_training_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_records_course FOREIGN KEY (course_id) REFERENCES training_courses(id) ON DELETE CASCADE,
  CONSTRAINT uq_training_records_employee_course UNIQUE (employee_id, course_id)
);
CREATE INDEX IF NOT EXISTS idx_training_records_deleted_at ON training_records (deleted_at);
CREATE INDEX IF NOT EXISTS idx_training_records_employee_id ON training_records (employee_id);
CREATE INDEX IF NOT EXISTS idx_training_records_course_id ON training_records (course_id);
CREATE INDEX IF NOT EXISTS idx_training_records_status ON training_records (status);
CREATE INDEX IF NOT EXISTS idx_training_records_start_date ON training_records (start_date);

-- ============================================================================
-- 初始数据插入
-- ============================================================================

-- 插入默认职位
INSERT INTO positions (code, name, description, department_id, level, min_salary, max_salary)
SELECT 'DEV_JUNIOR', '初级开发工程师', '初级软件开发工程师', d.id, 1, 8000.00, 12000.00
FROM departments d WHERE d.code = 'IT'
ON CONFLICT (code) DO NOTHING;

INSERT INTO positions (code, name, description, department_id, level, min_salary, max_salary)
SELECT 'DEV_SENIOR', '高级开发工程师', '高级软件开发工程师', d.id, 3, 15000.00, 25000.00
FROM departments d WHERE d.code = 'IT'
ON CONFLICT (code) DO NOTHING;

-- 插入默认薪资等级
INSERT INTO salary_grades (code, name, description, min_salary, max_salary, is_active) VALUES
('GRADE_1', '一级薪资', '初级员工薪资等级', 5000.00, 10000.00, TRUE),
('GRADE_2', '二级薪资', '中级员工薪资等级', 10000.00, 20000.00, TRUE),
('GRADE_3', '三级薪资', '高级员工薪资等级', 20000.00, 35000.00, TRUE),
('GRADE_4', '四级薪资', '专家级员工薪资等级', 35000.00, 50000.00, TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认请假类型
INSERT INTO leave_types (code, name, description, max_days_per_year, is_paid, requires_approval, is_active) VALUES
('ANNUAL', '年假', '带薪年假', 15, TRUE, TRUE, TRUE),
('SICK', '病假', '病假', 30, TRUE, TRUE, TRUE),
('PERSONAL', '事假', '个人事假', NULL, FALSE, TRUE, TRUE),
('MATERNITY', '产假', '产假', 128, TRUE, TRUE, TRUE),
('PATERNITY', '陪产假', '陪产假', 15, TRUE, TRUE, TRUE)
ON CONFLICT (code) DO NOTHING;

-- 插入默认培训课程
INSERT INTO training_courses (code, name, description, instructor, duration_hours, max_participants, cost, is_active) VALUES
('ORIENTATION', '新员工入职培训', '新员工入职培训课程', 'HR部门', 8, 20, 0.00, TRUE),
('SAFETY', '安全培训', '工作场所安全培训', '安全部门', 4, 50, 0.00, TRUE),
('TECH_BASIC', '技术基础培训', '基础技术技能培训', '技术部门', 16, 15, 500.00, TRUE)
ON CONFLICT (code) DO NOTHING;

-- ============================================================================
-- 结束
-- ============================================================================