#!/bin/bash
docker rm -f fafadoc
docker run --name fafadoc -d -p 8888:80 hunterhug/fafadoc