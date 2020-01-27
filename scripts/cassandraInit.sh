#!/bin/sh

sudo docker exec -it bapiCassandra /bin/mkdir /initdb -p

files=$(ls initdb)
for file in ${files}; do
	echo "--->executing $file"
	sudo docker cp initdb/${file} bapiCassandra:/initdb/${file}
	sudo docker exec -it bapiCassandra cqlsh -f /initdb/${file}
done
