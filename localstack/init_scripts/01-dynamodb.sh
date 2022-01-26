#!/bin/sh

awslocal dynamodb create-table --cli-input-json file://tmp/dynamo_table/access_log_table.json
awslocal dynamodb create-table --cli-input-json file://tmp/dynamo_table/api_routing_table.json
