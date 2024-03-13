-- Drop existing tables if they exist to avoid conflicts
DROP TABLE IF EXISTS employees CASCADE;
DROP TABLE IF EXISTS departments CASCADE;
DROP VIEW IF EXISTS department_summary;

-- Create departments table
CREATE TABLE departments (
    department_id SERIAL PRIMARY KEY,
    department_name VARCHAR(255) NOT NULL
);

-- Create employees table
CREATE TABLE employees (
    employee_id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    department_id INTEGER,
    email VARCHAR(255),
    CONSTRAINT fk_department
        FOREIGN KEY(department_id) 
	        REFERENCES departments(department_id)
);

-- Insert sample data into departments
INSERT INTO departments (department_name) VALUES
('Human Resources'),
('Research and Development'),
('Sales and Marketing');

-- Insert sample data into employees
INSERT INTO employees (first_name, last_name, department_id, email) VALUES
('John', 'Doe', 1, 'john.doe@example.com'),
('Jane', 'Doe', 2, 'jane.doe@example.com'),
('Jim', 'Beam', 3, 'jim.beam@example.com');

-- Create an index on the email column of the employees table
CREATE INDEX idx_employee_email ON employees(email);

-- Create a view to summarize department information
CREATE VIEW department_summary AS
SELECT d.department_name, COUNT(e.employee_id) AS number_of_employees
FROM departments d
LEFT JOIN employees e ON d.department_id = e.department_id
GROUP BY d.department_name;

-- Example of a simple SELECT to verify the view
SELECT * FROM department_summary;

