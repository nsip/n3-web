#!/bin/bash 

# shell script to load set of linked / useful demo data files

echo 'generating demo environment token'
curl -s -X POST -d 'userName=n3Demo' -d 'contextName=mySchool' \
localhost:1323/admin/newdemocontext > token.json


echo 'publishing school xapi data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/xapi/xapi_en.json \
http://localhost:1323/n3/publish

# echo 'publishing school naplan data...'
# time curl -X POST -H "Content-Type: application/json"  \
# -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
# -d @sample_data/naplan/sif.json \
# http://localhost:1323/n3/publish

echo 'publishing teaching subject data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/subjects/subjects.json \
http://localhost:1323/n3/publish

echo 'publishing lesson data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/lessons/lessons.json \
http://localhost:1323/n3/publish

echo 'publishing curriculum data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_history_stage1.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_history_stage2.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_history_stage3.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_history_stage4.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_history_stage5.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_geography_stage1.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_geography_stage2.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_geography_stage3.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_geography_stage4.json \
http://localhost:1323/n3/publish

time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/curriculum/hsie_geography_stage5.json \
http://localhost:1323/n3/publish

echo 'publishing school sif data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbCIsInVuYW1lIjoibjNEZW1vIn0.VTD8C6pwbkQ32u-vvuHnxq3xijdwNTd54JAyt1iLF3I" \
-d @sample_data/sif/sif.json \
http://localhost:1323/n3/publish





