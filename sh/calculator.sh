#!/bin/bash

protoc --proto_path calculator/proto --go_out=plugins=grpc:./calculator/pb calculator.proto
