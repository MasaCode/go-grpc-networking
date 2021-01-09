#!/bin/bash

protoc --proto_path greet/proto --go_out=plugins=grpc:./greet/pb greet.proto
