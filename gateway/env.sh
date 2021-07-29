#!/bin/bash
export READTIMEOUT=5
export READHEADERTIMEOUT=5
export WRITETIMEOUT=20
export IDLETIMEOUT=5
export MAXHEADERBYTES="1<<20"
export REDIS_HOST="localhost:6379"
export LOG_PATH="./log.csv"