DROP TABLE IF EXISTS bapi.user_created_events;
CREATE TABLE bapi.user_created_events (
	user_uuid UUID,
	event_uuid UUID,
	PRIMARY KEY (user_uuid, event_uuid)
);
