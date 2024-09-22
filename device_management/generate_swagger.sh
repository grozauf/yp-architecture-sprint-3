#!/bin/sh

mv main.go main.go.origin
oapi-codegen -generate gin -o ./swagger/api.go ./swagger.yaml
oapi-codegen -generate types -o ./swagger/models.go ./swagger.yaml
mv main.go.origin main.go
