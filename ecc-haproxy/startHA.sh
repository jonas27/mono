#!/bin/bash
docker build -t my-haproxy .
docker run -d --name my-running-haproxy my-haproxy -p 9090:9090