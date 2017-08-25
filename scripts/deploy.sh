#!/bin/sh

./scripts/build-web.sh && rsync -av . deploy@ta.do:/opt/sql-unit-test

