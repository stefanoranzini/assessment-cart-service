#!/usr/bin/env bash

# Useful to be sure to have always fresher code up&running. A real deployment scenario will have a docker image only with freshere binaries and no golang compiler dependencies.
./scripts/build.sh

./bin/migrations
./bin/cart