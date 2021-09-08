#!/bin/bash
export READTIMEOUT=5
export READHEADERTIMEOUT=5
export WRITETIMEOUT=20
export IDLETIMEOUT=5
export MAXHEADERBYTES="1<<20"
export API_DB_TYPE="REDIS"; \
export REDIS_HOST="localhost"; \
export REDIS_PORT="6379"; \
