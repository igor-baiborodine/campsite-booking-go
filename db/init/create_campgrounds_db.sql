CREATE DATABASE campgrounds;
CREATE USER campgrounds_user WITH ENCRYPTED PASSWORD 'campgrounds_pass';
ALTER DATABASE campgrounds OWNER TO campgrounds_user;
GRANT USAGE ON SCHEMA public TO campgrounds_user;

CREATE EXTENSION IF NOT EXISTS moddatetime;
