#!/bin/bash
sudo docker build -t josecarlosme/tungsteno-download-api . --no-cache
sudo docker push josecarlosme/tungsteno-download-api