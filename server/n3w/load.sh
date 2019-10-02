#!/bin/bash 

# shell script to load set of linked / useful demo data files

echo 'publishing school xapi data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
-d @sample_data/xapi/xapi_en.json \
http://localhost:1323/n3/publish

# echo 'publishing school naplan data...'
# time curl -X POST -H "Content-Type: application/json"  \
# -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
# -d @sample_data/naplan/sif.json \
# http://localhost:1323/n3/publish

echo 'publishing teaching subject data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
-d @sample_data/subjects/subjects.json \
http://localhost:1323/n3/publish

echo 'publishing lesson data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
-d @sample_data/lessons/lessons.json \
http://localhost:1323/n3/publish

echo 'publishing curriculum data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
-d @sample_data/curriculum/hsie_history_stage4.json \
http://localhost:1323/n3/publish

# time curl -X POST -H "Content-Type: application/json"  \
# -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
# -d @sample_data/curriculum/content.json \
# http://localhost:1323/n3/publish

echo 'publishing school sif data...'
time curl -X POST -H "Content-Type: application/json"  \
-H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJkZW1vIiwiY25hbWUiOiJteVNjaG9vbDIiLCJ1bmFtZSI6Im1hdHRmIn0.TY_C1zAFolhAhyv-uzxjiUw0eb5Wq32rqfdmxXuO7bM" \
-d @sample_data/sif/sif.json \
http://localhost:1323/n3/publish





