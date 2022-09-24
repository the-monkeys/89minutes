#!/bin/bash

echo "compiling proto files"
protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:pb

echo "proto files are compiled successfully"

read -p "do you want to run the test cases? (Y/N)" test

if (( $test == "y" || $test == "Y" )); then
    	echo "you want to test the application"
else
    	echo "starting the server"
fi


