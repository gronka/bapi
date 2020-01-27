DROP TABLE IF EXISTS bapi.events;
CREATE TABLE bapi.events (
	event_uuid UUID,
	end_time BIGINT,
	start_time BIGINT,
	tac_uuid UUID,
	tac_name TEXT,
	title TEXT,
	pic_url TEXT,
	address TEXT,
	lat DOUBLE,
	lng DOUBLE,
	tz_offset BIGINT,
	tz_id TEXT,
	tz_name TEXT,
	quick_info TEXT,
	long_info TEXT,
	phone TEXT,
	tier INT,
	rating BIGINT,
	admins SET<UUID>,
	organizers SET<UUID>,
	PRIMARY KEY (event_uuid)
);


DROP TABLE IF EXISTS bapi.event_has_admins;
CREATE TABLE bapi.event_has_admins (
	event_uuid UUID,
	user_uuid UUID,
	PRIMARY KEY (event_uuid, user_uuid)
);


// NOTE: sort by end time so that we pull events which have already started. Leave sorting by start_time to the client
DROP TABLE IF EXISTS bapi.user_admins_events;
CREATE TABLE bapi.user_admins_events (
	user_uuid UUID,
	end_time BIGINT,
	start_time BIGINT,
	event_uuid UUID,
	title TEXT,
	tac_name TEXT,
	address TEXT,
	quick_info TEXT,
	PRIMARY KEY (user_uuid, end_time, event_uuid)
) WITH CLUSTERING ORDER BY (end_time DESC);


DROP TABLE IF EXISTS bapi.user_rsvped_events;
CREATE TABLE bapi.user_rsvped_events (
	user_uuid UUID,
	end_time BIGINT,
	start_time BIGINT,
	event_uuid UUID,
	title TEXT,
	tac_name TEXT,
	address TEXT,
	quick_info TEXT,
	rsvp TEXT,
	PRIMARY KEY (user_uuid, end_time, event_uuid)
) WITH CLUSTERING ORDER BY (end_time DESC);


DROP TABLE IF EXISTS bapi.user_organizing_events;
CREATE TABLE bapi.user_organizing_events (
	user_uuid UUID,
	end_time BIGINT,
	start_time BIGINT,
	event_uuid UUID,
	title TEXT,
	tac_name TEXT,
	address TEXT,
	quick_info TEXT,
	PRIMARY KEY (user_uuid, end_time, event_uuid)
) WITH CLUSTERING ORDER BY (end_time DESC);


DROP TABLE IF EXISTS bapi.event_declined_users;
CREATE TABLE bapi.event_declined_users (
	event_uuid UUID,
	user_uuid UUID,
	PRIMARY KEY (event_uuid, user_uuid)
);
