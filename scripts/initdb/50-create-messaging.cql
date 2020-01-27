DROP TABLE IF EXISTS bapi.convos_by_time;
CREATE TABLE bapi.convos_by_time (
	user_uuid UUID,
	last_msg_time bigint,
	convo_uuid UUID,
	PRIMARY KEY (user_uuid, last_msg_time, convo_uuid)
) WITH CLUSTERING ORDER BY (last_msg_time DESC);


/*admins will be able to see real_uuid on screen*/
DROP TABLE IF EXISTS bapi.convo_msgs;
CREATE TABLE bapi.convo_msgs (
	convo_uuid UUID,
	msg_id int,
	time_sent bigint,
	msg_uuid UUID,
	apparent_uuid UUID,
	real_uuid UUID,
	body text,
	PRIMARY KEY (convo_uuid, msg_id)
) WITH CLUSTERING ORDER BY (msg_id DESC);


DROP TABLE IF EXISTS bapi.convo_msg_id_counter;
CREATE TABLE bapi.convo_msg_id_counter (
	convo_uuid UUID,
	msg_id counter,
	PRIMARY KEY (convo_uuid)
);


/*We can create a unique hash for each group conversation by sorting the */
/*user_uuids then hashing.*/
DROP TABLE IF EXISTS bapi.participant_hash_to_convo_uuid;
CREATE TABLE bapi.participant_hash_to_convo_uuid (
	participant_hash blob,
	convo_uuid UUID,
	PRIMARY KEY (participant_hash)
);


DROP TABLE IF EXISTS bapi.convo_uuid_to_user_uuid;
CREATE TABLE bapi.convo_uuid_to_user_uuid (
	convo_uuid UUID,
	user_uuid UUID,
	role text,
	status text,
	PRIMARY KEY (convo_uuid, user_uuid)
);

