DROP TABLE IF EXISTS bapi.user_has_mapsac_token;
CREATE TABLE bapi.user_has_mapsac_token (
	user_uuid uuid,
	timeout bigint,
	session_token uuid,
	consumed boolean,
	PRIMARY KEY (user_uuid, timeout, session_token)
) WITH CLUSTERING ORDER BY (timeout DESC);


DROP TABLE IF EXISTS bapi.user_has_tacs;
CREATE TABLE bapi.user_has_tacs (
	user_uuid uuid,
	tac_uuid uuid,
	last_used bigint,
	name text,
	address text,
	lat double,
	lng double,
	tz_offset int,
	tz_id text,
	PRIMARY KEY (user_uuid, tac_uuid)
);


DROP TABLE IF EXISTS bapi.tac_has_gg_place;
CREATE TABLE bapi.tac_has_gg_place (
	tac_uuid uuid,
	gg_place_id text,
	gg_lat double,
	gg_lng double,
	gg_name text,
	gg_formatted_address text,
	maps_url text,
	tz_offset int,
	tz_id text,
	PRIMARY KEY (tac_uuid)
);
CREATE INDEX ON bapi.tac_has_gg_place(gg_place_id);


UPDATE bapi.user_has_tacs SET 
	last_used = 1575317500996, 
	name = 'The Digs',
	address = '3721 Kirkwood Hwy, Wilmington, DE 19808, USA',
	lat = 39.73315,
	lng = -75.64088,
	tz_offset = -300
WHERE user_uuid=22222222-2222-2222-2222-222222222222 AND tac_uuid=22222222-2222-2222-2222-222222333333;


UPDATE bapi.user_has_tacs SET 
	last_used = 1579205795626, 
	name = 'Christian Fellowship',
	address = '1501 Athens Dr, Raleigh, NC 27606, USA',
	lat = 35.76992,
	lng = -78.71026,
	tz_offset = -300
WHERE user_uuid=22222222-2222-2222-2222-222222222222 AND tac_uuid=22222222-2222-2222-2222-333222333333;


UPDATE bapi.user_has_tacs SET 
	last_used = 1579205758164, 
	name = 'Baptist Church',
	address = '1601 Athens Dr, Raleigh, NC 27606, USA',
	lat = 35.76992,
	lng = -78.71026,
	tz_offset = -300
WHERE user_uuid=22222222-2222-2222-2222-222222222222 AND tac_uuid=22222222-2222-2222-2222-444222333333;


UPDATE bapi.user_has_tacs SET 
	last_used = 1579205823741, 
	name = 'Community Garden',
	address = '1321 Athens Dr, Raleigh, NC 27606, USA',
	lat = 35.77119,
	lng = -78.71221,
	tz_offset = -300
WHERE user_uuid=22222222-2222-2222-2222-222222222222 AND tac_uuid=22222222-2222-2222-2222-555222333333;
