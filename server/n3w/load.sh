#!/bin/bash 

# shell script to load set of linked / useful demo data files

echo 'publishing school xapi data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/xapi/xapi.json \
http://localhost:1323/publish

echo 'publishing school naplan data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/naplan/sif.json \
http://localhost:1323/publish

echo 'publishing teaching subject data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/subjects/subjects.json \
http://localhost:1323/publish

echo 'publishing lesson data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/lessons/lessons.json \
http://localhost:1323/publish

echo 'publishing curriculum data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/curriculum/overview.json \
http://localhost:1323/publish

time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/curriculum/content.json \
http://localhost:1323/publish

echo 'publishing school sif data...'
time curl -X POST -H "Content-Type: application/json"  \
-d @sample_data/sif/sif.json \
http://localhost:1323/publish





