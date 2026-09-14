package models

const CreateTablesSQL = `

CREATE TABLE IF NOT EXISTS states (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,
	state_abbreviation TEXT NOT NULL,
	CONSTRAINT pk_states PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS cities (
	id INTEGER NOT NULL,
	name TEXT NOT NULL,
	state_id INTEGER NOT NULL,
	CONSTRAINT pk_cities PRIMARY KEY (id),
	CONSTRAINT fk_cities_state_id FOREIGN KEY (state_id) REFERENCES states (id)
);

CREATE TABLE IF NOT EXISTS users (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	cpf TEXT NOT NULL,
	CONSTRAINT pk_users PRIMARY KEY (id),
	CONSTRAINT uq_users_email UNIQUE (email),
	CONSTRAINT uq_users_cpf UNIQUE (cpf)
);

CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	user_id INTEGER NOT NULL,
	login_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expired_at TIMESTAMPTZ,
	logout_at TIMESTAMPTZ,
	CONSTRAINT pk_sessions PRIMARY KEY (id),
	CONSTRAINT fk_sessions_user_id FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS pickup_location (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	city_id INTEGER NOT NULL,
	state_id INTEGER NOT NULL,
	district TEXT NOT NULL,
	street TEXT NOT NULL,
	building_number TEXT NOT NULL,
	CONSTRAINT pk_pickup_location PRIMARY KEY (id),
	CONSTRAINT fk_pickup_location_city_id FOREIGN KEY (city_id) REFERENCES cities (id),
	CONSTRAINT fk_pickup_location_state_id FOREIGN KEY (state_id) REFERENCES states (id)
);

CREATE TABLE IF NOT EXISTS vehicles (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	brand TEXT NOT NULL,
	model TEXT NOT NULL,
	chassi TEXT NOT NULL,
	year INTEGER NOT NULL,
	color TEXT NOT NULL,
	doors_amount INTEGER NOT NULL,
	seats_amount INTEGER NOT NULL,
	has_air_conditioning BOOLEAN NOT NULL DEFAULT FALSE,
	pickup_location_id INTEGER NOT NULL,
	CONSTRAINT pk_vehicles PRIMARY KEY (id),
	CONSTRAINT uq_vehicles_chassi UNIQUE (chassi),
	CONSTRAINT fk_vehicles_pickup_location_id FOREIGN KEY (pickup_location_id) REFERENCES pickup_location (id)
);

CREATE TABLE IF NOT EXISTS allocations (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	renter_id INTEGER NOT NULL,
	vehicle_id INTEGER NOT NULL,
	pick_up_date TIMESTAMP NOT NULL,
	devolution_date TIMESTAMP,
	estimated_evolution_date TIMESTAMP,
	CONSTRAINT pk_allocations PRIMARY KEY (id),
	CONSTRAINT fk_allocations_renter_id FOREIGN KEY (renter_id) REFERENCES users (id),
	CONSTRAINT fk_allocations_vehicle_id FOREIGN KEY (vehicle_id) REFERENCES vehicles (id)
);

CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER GENERATED ALWAYS AS IDENTITY,
	title VARCHAR(100) NOT NULL,
	description TEXT,
	status BOOLEAN NOT NULL DEFAULT FALSE,
	CONSTRAINT pk_tasks PRIMARY KEY (id)
);
`