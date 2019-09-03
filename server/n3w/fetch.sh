#!/bin/bash

# simple test of finder on n3w

echo 'running staff traversal query...'
time curl -s -X GET -H "Content-Type: application/json"  \
http://localhost:1323/staffTraversal?staffid\=A4F0069E-D3B8-4822-BDD9-4D649E2A47FD \
| jq . 
# > results.json

