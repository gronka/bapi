/* start_time is a combo of weekday + starting time of day */
DROP TABLE IF EXISTS bapi.planner_updates;
CREATE TABLE bapi.planner_updates (
	owner_uuid UUID,
	available_updates counter,
	filled_updates counter,
	PRIMARY KEY (owner_uuid)
);


DROP TABLE IF EXISTS bapi.planner_available_times;
CREATE TABLE bapi.planner_available_times (
	owner_uuid UUID,
	weekday int,
	start_mm int,
	end_mm int,
	tz_offset int,
	PRIMARY KEY (owner_uuid, weekday, start_mm)
);
/*CREATE INDEX ON bapi.planner_available_times(owner_uuid);*/


DROP TABLE IF EXISTS bapi.planner_filled_times;
CREATE TABLE bapi.planner_filled_times (
	owner_uuid UUID,
	day_of_year int,
	start_mm int,
	status int,
	appt_uuid UUID,
	end_mm int,
	reason text,
	tz_offset int,
	requestee_uuid UUID,
	canceled_by text,
	PRIMARY KEY (owner_uuid, day_of_year, start_mm, status, appt_uuid)
);
/*CREATE INDEX ON bapi.planner_filled_times(appt_uuid);*/
