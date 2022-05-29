go test -coverprofile=tm_coverage.out ./app/trademgr
go tool cover -func=tm_coverage.out > unittests-coverage.txt
go tool cover -html=tm_coverage.out -o unittests-coverage.html
rm tm_coverage.out
