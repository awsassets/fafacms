#!/bin/bash
docker run -it --rm -p 8888:3000 --name=docsify -v $(pwd):/docs hunterhug/docsify