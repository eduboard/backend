#!/usr/bin/env bash
docker build . -t eduboard-backend:`git rev-parse HEAD`
