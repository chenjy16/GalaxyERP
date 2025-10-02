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
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  is_active BOOLEAN DEFAULT TRUE,
  code VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT NULL,
  department_id INTEGER NOT NULL,
  CONSTRAINT uq_positions_code UNIQUE (code),
  CONSTRAINT fk_positions_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE CASCADE,
  CONSTRAINT fk_positions_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_positions_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_positions_deleted_at ON positions (deleted_at);
CREATE INDEX IF NOT EXISTS idx_positions_is_active ON positions (is_active);
CREATE INDEX IF NOT EXISTS idx_positions_department_id ON positions (department_id);
CREATE INDEX IF NOT EXISTS idx_positions_code ON positions (code);

-- ============================================================================
-- 员工管理
-- ============================================================================

-- 员工表
CREATE TABLE IF NOT EXISTS employees (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  code VARCHAR(50) NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  phone VARCHAR(20) NULL,
  date_of_birth TIMESTAMP WITH TIME ZONE NULL,
  gender VARCHAR(10) NULL,
  hire_date TIMESTAMP WITH TIME ZONE NULL,
  department_id INTEGER NULL,
  position_id INTEGER NULL,
  manager_id INTEGER NULL,
  status VARCHAR(50) DEFAULT 'active',
  emergency_contact VARCHAR(255) NULL,
  id_number VARCHAR(100) NULL,
  address TEXT NULL,
  CONSTRAINT uq_employees_code UNIQUE (code),
  CONSTRAINT uq_employees_email UNIQUE (email),
  CONSTRAINT fk_employees_department FOREIGN KEY (department_id) REFERENCES departments(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_manager FOREIGN KEY (manager_id) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employees_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_employees_deleted_at ON employees (deleted_at);
CREATE INDEX IF NOT EXISTS idx_employees_code ON employees (code);
CREATE INDEX IF NOT EXISTS idx_employees_email ON employees (email);
CREATE INDEX IF NOT EXISTS idx_employees_department_id ON employees (department_id);
CREATE INDEX IF NOT EXISTS idx_employees_position_id ON employees (position_id);
CREATE INDEX IF NOT EXISTS idx_employees_manager_id ON employees (manager_id);
CREATE INDEX IF NOT EXISTS idx_employees_status ON employees (status);

-- ============================================================================
-- 薪资管理
-- ============================================================================



-- ============================================================================
-- 考勤管理
-- ============================================================================

-- 考勤记录表
CREATE TABLE IF NOT EXISTS attendances (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  check_in_time TIMESTAMP WITH TIME ZONE NULL,
  check_out_time TIMESTAMP WITH TIME ZONE NULL,
  status VARCHAR(50) DEFAULT 'present',
  notes TEXT NULL,
  CONSTRAINT fk_attendances_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_attendances_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_attendances_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_attendances_deleted_at ON attendances (deleted_at);
CREATE INDEX IF NOT EXISTS idx_attendances_employee_id ON attendances (employee_id);
CREATE INDEX IF NOT EXISTS idx_attendances_date ON attendances (date);
CREATE INDEX IF NOT EXISTS idx_attendances_status ON attendances (status);

-- ============================================================================
-- 请假管理
-- ============================================================================

-- 请假申请表
CREATE TABLE IF NOT EXISTS leaves (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  leave_type VARCHAR(50) NOT NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  days DECIMAL(5,2) NOT NULL,
  reason TEXT NOT NULL,
  status VARCHAR(50) DEFAULT 'pending',
  approved_by INTEGER NULL,
  approved_at TIMESTAMP WITH TIME ZONE NULL,
  CONSTRAINT fk_leaves_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_leaves_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_leaves_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_leaves_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_leaves_deleted_at ON leaves (deleted_at);
CREATE INDEX IF NOT EXISTS idx_leaves_employee_id ON leaves (employee_id);
CREATE INDEX IF NOT EXISTS idx_leaves_status ON leaves (status);
CREATE INDEX IF NOT EXISTS idx_leaves_start_date ON leaves (start_date);
CREATE INDEX IF NOT EXISTS idx_leaves_end_date ON leaves (end_date);

-- ============================================================================
-- 绩效管理
-- ============================================================================

-- 绩效评估表
CREATE TABLE IF NOT EXISTS performance_reviews (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  reviewer_id INTEGER NOT NULL,
  review_date TIMESTAMP WITH TIME ZONE NOT NULL,
  review_period VARCHAR(100) NOT NULL,
  score DECIMAL(5,2) NOT NULL,
  comments TEXT NULL,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_performance_reviews_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_reviewer FOREIGN KEY (reviewer_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_reviews_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_performance_reviews_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_deleted_at ON performance_reviews (deleted_at);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_employee_id ON performance_reviews (employee_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_reviewer_id ON performance_reviews (reviewer_id);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_status ON performance_reviews (status);
CREATE INDEX IF NOT EXISTS idx_performance_reviews_review_date ON performance_reviews (review_date);

-- ============================================================================
-- 培训管理
-- ============================================================================

-- 技能表
CREATE TABLE IF NOT EXISTS skills (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  name VARCHAR(100) NOT NULL,
  description TEXT NULL,
  category VARCHAR(50) NULL,
  CONSTRAINT fk_skills_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_skills_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_skills_deleted_at ON skills (deleted_at);
CREATE INDEX IF NOT EXISTS idx_skills_name ON skills (name);
CREATE INDEX IF NOT EXISTS idx_skills_category ON skills (category);

-- 员工技能表
CREATE TABLE IF NOT EXISTS employee_skills (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  skill_id INTEGER NOT NULL,
  level VARCHAR(50) NOT NULL,
  acquired_date TIMESTAMP WITH TIME ZONE NULL,
  CONSTRAINT fk_employee_skills_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_employee_skills_skill FOREIGN KEY (skill_id) REFERENCES skills(id) ON DELETE CASCADE,
  CONSTRAINT fk_employee_skills_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_employee_skills_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  UNIQUE(employee_id, skill_id)
);
CREATE INDEX IF NOT EXISTS idx_employee_skills_deleted_at ON employee_skills (deleted_at);
CREATE INDEX IF NOT EXISTS idx_employee_skills_employee_id ON employee_skills (employee_id);
CREATE INDEX IF NOT EXISTS idx_employee_skills_skill_id ON employee_skills (skill_id);
CREATE INDEX IF NOT EXISTS idx_employee_skills_level ON employee_skills (level);

-- 工资单表
CREATE TABLE IF NOT EXISTS payrolls (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  pay_period TIMESTAMP WITH TIME ZONE NOT NULL,
  basic_salary DECIMAL(15,2) NOT NULL DEFAULT 0,
  overtime_pay DECIMAL(15,2) NOT NULL DEFAULT 0,
  bonus DECIMAL(15,2) NOT NULL DEFAULT 0,
  deductions DECIMAL(15,2) NOT NULL DEFAULT 0,
  net_pay DECIMAL(15,2) NOT NULL DEFAULT 0,
  status VARCHAR(50) DEFAULT 'draft',
  CONSTRAINT fk_payrolls_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_payrolls_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_payrolls_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_payrolls_deleted_at ON payrolls (deleted_at);
CREATE INDEX IF NOT EXISTS idx_payrolls_employee_id ON payrolls (employee_id);
CREATE INDEX IF NOT EXISTS idx_payrolls_pay_period ON payrolls (pay_period);
CREATE INDEX IF NOT EXISTS idx_payrolls_status ON payrolls (status);

-- 培训表
CREATE TABLE IF NOT EXISTS trainings (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NULL,
  trainer VARCHAR(255) NULL,
  start_date TIMESTAMP WITH TIME ZONE NOT NULL,
  end_date TIMESTAMP WITH TIME ZONE NOT NULL,
  location VARCHAR(255) NULL,
  max_participants INTEGER NULL,
  status VARCHAR(50) DEFAULT 'planned',
  CONSTRAINT fk_trainings_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_trainings_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_trainings_deleted_at ON trainings (deleted_at);
CREATE INDEX IF NOT EXISTS idx_trainings_status ON trainings (status);
CREATE INDEX IF NOT EXISTS idx_trainings_start_date ON trainings (start_date);
CREATE INDEX IF NOT EXISTS idx_trainings_end_date ON trainings (end_date);

-- 培训参与者表
CREATE TABLE IF NOT EXISTS training_participants (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  training_id INTEGER NOT NULL,
  employee_id INTEGER NOT NULL,
  registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  attendance_status VARCHAR(50) DEFAULT 'registered',
  completion_status VARCHAR(50) DEFAULT 'not_started',
  score DECIMAL(5,2) NULL,
  feedback TEXT NULL,
  CONSTRAINT fk_training_participants_training FOREIGN KEY (training_id) REFERENCES trainings(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_participants_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_training_participants_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_training_participants_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL,
  UNIQUE(training_id, employee_id)
);
CREATE INDEX IF NOT EXISTS idx_training_participants_deleted_at ON training_participants (deleted_at);
CREATE INDEX IF NOT EXISTS idx_training_participants_training_id ON training_participants (training_id);
CREATE INDEX IF NOT EXISTS idx_training_participants_employee_id ON training_participants (employee_id);
CREATE INDEX IF NOT EXISTS idx_training_participants_attendance ON training_participants (attendance_status);
CREATE INDEX IF NOT EXISTS idx_training_participants_completion ON training_participants (completion_status);

-- 加班记录表
CREATE TABLE IF NOT EXISTS overtime_records (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  date TIMESTAMP WITH TIME ZONE NOT NULL,
  start_time TIMESTAMP WITH TIME ZONE NOT NULL,
  end_time TIMESTAMP WITH TIME ZONE NOT NULL,
  hours DECIMAL(5,2) NOT NULL,
  reason TEXT NULL,
  status VARCHAR(50) DEFAULT 'pending',
  approved_by INTEGER NULL,
  approved_at TIMESTAMP WITH TIME ZONE NULL,
  CONSTRAINT fk_overtime_records_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_overtime_records_approved_by FOREIGN KEY (approved_by) REFERENCES employees(id) ON DELETE SET NULL,
  CONSTRAINT fk_overtime_records_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_overtime_records_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_overtime_records_deleted_at ON overtime_records (deleted_at);
CREATE INDEX IF NOT EXISTS idx_overtime_records_employee_id ON overtime_records (employee_id);
CREATE INDEX IF NOT EXISTS idx_overtime_records_date ON overtime_records (date);
CREATE INDEX IF NOT EXISTS idx_overtime_records_status ON overtime_records (status);

-- 绩效目标表
CREATE TABLE IF NOT EXISTS performance_goals (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE NULL,
  created_by INTEGER NULL,
  updated_by INTEGER NULL,
  employee_id INTEGER NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NULL,
  target_date TIMESTAMP WITH TIME ZONE NOT NULL,
  weight DECIMAL(5,2) DEFAULT 1.0,
  status VARCHAR(50) DEFAULT 'active',
  achievement_rate DECIMAL(5,2) DEFAULT 0.0,
  CONSTRAINT fk_performance_goals_employee FOREIGN KEY (employee_id) REFERENCES employees(id) ON DELETE CASCADE,
  CONSTRAINT fk_performance_goals_created_by FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL,
  CONSTRAINT fk_performance_goals_updated_by FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS idx_performance_goals_deleted_at ON performance_goals (deleted_at);
CREATE INDEX IF NOT EXISTS idx_performance_goals_employee_id ON performance_goals (employee_id);
CREATE INDEX IF NOT EXISTS idx_performance_goals_status ON performance_goals (status);
CREATE INDEX IF NOT EXISTS idx_performance_goals_target_date ON performance_goals (target_date);

-- ============================================================================
-- 初始数据
-- ============================================================================

-- 职位数据
INSERT INTO positions (code, name, description, department_id, created_at, updated_at) VALUES
('CEO', '首席执行官', '公司最高管理者', 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('CTO', '首席技术官', '技术部门负责人', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('DEV_SENIOR', '高级开发工程师', '负责核心系统开发', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('DEV_JUNIOR', '初级开发工程师', '参与系统开发工作', 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('HR_MANAGER', '人力资源经理', '人力资源部门负责人', 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('SALES_MANAGER', '销售经理', '销售部门负责人', 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 技能数据
INSERT INTO skills (name, description, category, created_at, updated_at) VALUES
('Go语言', 'Go编程语言开发技能', '编程语言', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('React', 'React前端框架开发技能', '前端技术', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('PostgreSQL', 'PostgreSQL数据库管理技能', '数据库', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('项目管理', '项目管理和团队协调能力', '管理技能', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('沟通协调', '团队沟通和协调能力', '软技能', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- 培训数据
INSERT INTO trainings (title, description, trainer, start_date, end_date, location, max_participants, status, created_at, updated_at) VALUES
('Go语言进阶培训', 'Go语言高级特性和最佳实践', '张三', '2024-03-01 09:00:00+08:00', '2024-03-05 17:00:00+08:00', '培训室A', 20, 'planned', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('React开发实战', 'React框架实战开发', '李四', '2024-03-15 09:00:00+08:00', '2024-03-19 17:00:00+08:00', '培训室B', 15, 'planned', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
('团队管理技能', '团队建设和管理技巧', '王五', '2024-04-01 09:00:00+08:00', '2024-04-03 17:00:00+08:00', '会议室C', 25, 'planned', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- ============================================================================
-- 结束
-- ============================================================================