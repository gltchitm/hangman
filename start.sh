#!/bin/sh

cd svelte
npm run build
cd ..
mix run --no-halt
