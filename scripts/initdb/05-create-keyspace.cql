// DROP KEYSPACE bapi;
CREATE KEYSPACE IF NOT EXISTS "bapi" WITH replication = { 'class': 'SimpleStrategy',  'replication_factor': '1' };
