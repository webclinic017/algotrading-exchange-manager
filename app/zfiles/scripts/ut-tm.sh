cd app/trademgr                 
wait
go test -coverprofile=coverage.out    
wait
go tool cover -html=coverage.out
wait