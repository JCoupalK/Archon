-- Create a table
CREATE TABLE test_table (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

-- Insert data into the table
INSERT INTO test_table (name) VALUES ('Alice');
INSERT INTO test_table (name) VALUES ('Bob');
INSERT INTO test_table (name) VALUES ('Jordan');
INSERT INTO test_table (name) VALUES ('Kevin');