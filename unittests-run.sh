go test -coverprofile=tm_coverage.out ./app/trademgr
go tool cover -func=tm_coverage.out > ./docs/ut-reports/ut-coverage-trademgr.txt
go tool cover -html=tm_coverage.out -o ./docs/ut-reports/ut-coverage-trademgr.html
rm tm_coverage.out
