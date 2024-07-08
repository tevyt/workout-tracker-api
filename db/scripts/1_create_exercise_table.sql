CREATE TABLE IF NOT EXISTS exercise (
    id SERIAL PRIMARY KEY,
    exercise_name VARCHAR(50) NOT NULL,
    increment SMALLINT NOT NULL
)