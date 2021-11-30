#!/bin/bash
for i in `seq 1 60`
do
    ab -c 100 -n 100  localhost:8081/postserver/get_post &
    sleep 1
done
