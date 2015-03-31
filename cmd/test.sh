#!/bin/bash

go build -o latest && \
    ./latest -repo=ghr -owner=tcnksm -new -debug 0.4.5

echo $? 
rm ./latest
