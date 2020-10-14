#!/bin/bash
./wait-for-it.sh postgres-app:5432 --timeout=30 && ./test
