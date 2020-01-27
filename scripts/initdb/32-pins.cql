DROP TABLE IF EXISTS bapi.event_pins;
CREATE TABLE bapi.event_pins (
	muid UUID,
	event_uuid UUID,
	end_time BIGINT,
	start_time BIGINT,
	tac_name TEXT,
	title TEXT,
	lat DOUBLE,
	lng DOUBLE,
	tz_offset BIGINT,
	rating BIGINT,
	PRIMARY KEY (muid)
);


DROP TABLE IF EXISTS bapi.event_has_muids;
CREATE TABLE bapi.event_has_muids (
	event_uuid UUID,
	muids SET<UUID>,
	muidBins SET<UUID>,
	PRIMARY KEY (event_uuid)
);


DROP TABLE IF EXISTS bapi.muid_counters;
CREATE TABLE bapi.muid_counters (
	muid_bin UUID,
	bronze_count COUNTER,
	silver_count COUNTER,
	gold_count COUNTER,
	diamond_count COUNTER,
	PRIMARY KEY (muid_bin)
);
