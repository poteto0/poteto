#!bin/bash
echo "start test"
go test ./... -cover -coverprofile=cover

echo "report"
go tool cover -func=cover

echo "post test"
rm -f cover.out