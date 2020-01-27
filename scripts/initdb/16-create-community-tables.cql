DROP TABLE IF EXISTS bapi.user_locs;
CREATE TABLE bapi.user_locs (
	geohash text,
	user_uuid UUID,
	lat float,
	lng float,
	PRIMARY KEY (geohash, user_uuid)
);
CREATE INDEX ON bapi.user_locs(user_uuid);


DROP TABLE IF EXISTS bapi.shouts;
CREATE TABLE bapi.shouts (
	geohash text,
	shout_uuid UUID,
	user_uuid UUID,
	datetime int,
	shout_text_num int,
	lat float,
	lng float,
	PRIMARY KEY (geohash, shout_uuid)
);
/*CREATE INDEX ON bapi.shouts(msg_uuid);*/


DROP TABLE IF EXISTS bapi.shout_locs;
CREATE TABLE bapi.shout_locs (
	geohash text,
	shout_uuid UUID,
	lat float,
	lng float,
	PRIMARY KEY (geohash, shout_uuid)
);


DROP TABLE IF EXISTS bapi.courses;
CREATE TABLE bapi.courses (
	geohash text,
	course_uuid UUID,
	lat float,
	lng float,

	about text,
	rules text,
	PRIMARY KEY (geohash, course_uuid)
);


DROP TABLE IF EXISTS bapi.course_locs;
CREATE TABLE bapi.course_locs (
	geohash text,
	course_uuid UUID,
	lat float,
	lng float,
	PRIMARY KEY (geohash, course_uuid)
);
CREATE INDEX ON bapi.course_locs(course_uuid);
