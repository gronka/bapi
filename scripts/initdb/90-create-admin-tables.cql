DROP TABLE IF EXISTS bapi.trainer_applications;
CREATE TABLE bapi.trainer_applications (
	user_uuid UUID,
	fullname text,
	pdga_id text,
	PRIMARY KEY (user_uuid)
);
