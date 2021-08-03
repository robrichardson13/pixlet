#!/bin/sh

echo "Building..."

cd render
go build
cd ../runtime
go build
cd ../encode
go build
cd ..
go build

echo "Build complete!"