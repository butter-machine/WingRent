CREATE TABLE IF NOT EXISTS plane_model(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS planes(
    id SERIAL PRIMARY KEY,
    plane_model_id INT NOT NULL,
    tail_number VARCHAR(10) NOT NULL UNIQUE,
    CONSTRAINT fk_plane_model
    FOREIGN KEY(plane_model_id)
    REFERENCES plane_model(id)
);
