#!/bin/sh
echo Start Rollout of $1 to $2
docker stop $2
docker rm $2
docker rmi $2
docker load -i $2.tar
docker run -p 80:9090 -d --name $2 $1
echo Rollout finished
