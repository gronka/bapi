DROP TABLE IF EXISTS bapi.user_has_stripe_id;
CREATE TABLE bapi.user_has_stripe_id (
	user_uuid UUID,
	stripe_id text,
	PRIMARY KEY (user_uuid)
);


DROP TABLE IF EXISTS bapi.user_has_orders;
CREATE TABLE bapi.user_has_orders (
	user_uuid UUID,
	order_uuid UUID,
	PRIMARY KEY (user_uuid, order_uuid)
);


DROP TABLE IF EXISTS bapi.order_has_payment_source;
CREATE TABLE bapi.order_has_payment_source (
	order_uuid UUID,
	source_id UUID,
	source_type text,
	PRIMARY KEY (order_uuid)
);


/*The order object should be able to capture items purchased or sessions logged*/
DROP TABLE IF EXISTS bapi.orders;
CREATE TABLE bapi.orders (
	order_uuid UUID,
	order_type text,  // appointment, purchase, advertising
	seller_uuid UUID,
	seller_type text,  // trainer, corporate account, etc
	description text,
	rate float,
	timespan_hours float,
	total_cost float,
	user_did_tip boolean,
	tip_amount float,
	/*item_uuid_list text,*/
	PRIMARY KEY (order_uuid, order_type)
);


DROP TABLE IF EXISTS bapi.checkout_sessions_owned_by_stripe_id;
CREATE TABLE bapi.checkout_sessions_owned_by_stripe_id (
	checkout_session_id text,
	stripe_id text,
	expire_time bigint,
	PRIMARY KEY (checkout_session_id, stripe_id)
);


DROP TABLE IF EXISTS bapi.stripe_id_has_destitute_payment_methods;
CREATE TABLE bapi.stripe_id_has_destitute_payment_methods (
	stripe_id text,
	payment_method_id text,
	brand text,
	exp_month int,
	exp_year int,
	last4 text,
	PRIMARY KEY (stripe_id, payment_method_id)
);


/*possible future option which could let us get all payment info in one query:*/
/*DROP TABLE IF EXISTS bapi.useruuid_to_payment_id;*/
/*CREATE TABLE bapi.useruuid_to_payment_id (*/
	/*user_uuid UUID,*/
	/*id string,*/
	/*for_vendor string, < would be stripe/paypal/etc*/
	/*PRIMARY KEY (user_uuid, id, for_vendor)*/
/*);*/
