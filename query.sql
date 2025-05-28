CREATE DATABASE iot;
USE iot;

CREATE TABLE pzem (
	voltage DOUBLE,
	current DOUBLE,
	power DOUBLE,
	energy DOUBLE,
	frequency DOUBLE,
	power_factor DOUBLE,
	timestamp DATETIME
);

CREATE TABLE xymd (
	temperature DOUBLE,
	humidity DOUBLE,
	timestamp DATETIME
);