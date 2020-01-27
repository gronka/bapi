#!/bin/sh
sudo docker build -t cassandrabop -f ./cassandra.dockerfile .
sudo docker run --name bapiCassandra \
	-p 9042:9042 \
	-v /docker/bapicassandra:/var/lib \
	cassandrabop

	#cassandra:3.11.4
	#-v /docker/cassandra/etc:/etc/cassandra \
	#--mount type=bind,source=$(pwd)/cassandra.yaml,target=/etc/cassandra/cassandra.yaml \
