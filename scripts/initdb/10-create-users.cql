DROP TABLE IF EXISTS bapi.users;
CREATE TABLE bapi.users (
	user_uuid UUID PRIMARY KEY,
	language text,
	password text,

	email text,
	email_verified boolean,
	calling_code text,
	iso_phone text,
	naked_phone text,
	phone_country_iso text,
	phone_verified boolean,

	fullname text,
	first_name text,
	middle_name text,
	last_name text,
	status text,
	about text,
	hometown text,
	profile_pic_url text,

	popularity float,
	lat float,
	lng float,
	geohash text,

	account_status text,
	last_online bigint,
	online_status text,
);
CREATE INDEX ON bapi.users(iso_phone);
CREATE INDEX ON bapi.users(email);
CREATE INDEX ON bapi.users(first_name);
CREATE INDEX ON bapi.users(middle_name);
CREATE INDEX ON bapi.users(last_name);


DROP TABLE IF EXISTS bapi.phone_verify;

CREATE TABLE bapi.phone_verify (
	user_uuid UUID PRIMARY KEY,
	installation_id text,
	code text,
	expire_time bigint,
	attempts int,
);
CREATE INDEX ON bapi.phone_verify(installation_id);
CREATE INDEX ON bapi.phone_verify(expire_time);


DROP TABLE IF EXISTS bapi.phone_verify_spam;

CREATE TABLE bapi.phone_verify_spam (
	installation_id text PRIMARY KEY,
	attempts int,
	last_update int,
);


INSERT INTO bapi.users (user_uuid) VALUES (11111111-1111-1111-1111-111111111111);
UPDATE bapi.users SET 
	phone_verified = true,
	calling_code = '1',
	iso_phone = 'us1',
	naked_phone = '1',
	phone_country_iso = 'us',
	account_status = 'FINE',
	fullname = 'notset',
	about = 'I am the first user here',
	password = '1'
WHERE user_uuid=11111111-1111-1111-1111-111111111111;

INSERT INTO bapi.users (user_uuid) VALUES (22222222-2222-2222-2222-222222222222);
UPDATE bapi.users SET 
	phone_verified = true,
	calling_code = '1',
	iso_phone = 'us1',
	naked_phone = '1',
	phone_country_iso = 'us',
	account_status = 'FINE',
	fullname = 'notset',
	about = 'I made my account second',
	password = '1'
WHERE user_uuid=22222222-2222-2222-2222-222222222222;

INSERT INTO bapi.users (user_uuid) VALUES (99999999-9999-9999-9999-999999999999);
UPDATE bapi.users SET 
	phone_verified = true,
	calling_code = '1',
	iso_phone = 'us9',
	naked_phone = '9',
	phone_country_iso = 'us',
	account_status = 'FINE',
	fullname = 'Terry Hacker',
	about = 'I made my account second',
	password = '9'
WHERE user_uuid=99999999-9999-9999-9999-999999999999;
