#!/usr/bin/env sh

set -e

npm --prefix svelte ci
npm --prefix svelte run build

cd server
go run .
cd ..
