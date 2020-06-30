#!/bin/sh

docker stop grafana-service && docker rm grafana-service

docker run -p 8080:8080 --env-file=./env --name grafana-service keptnsandbox/grafana-service:dev
