# UpDoc Guide

##

Test command

` go run ./cmd/server &`

`curl http://localhost:8080/v1/docs/622593`

## db

docker compose up -d db

docker exec -it $(docker compose ps -q db) \
 psql -U updoc updoc -c "select flag_id,page_id,status from flags;"

docker exec -it $(docker compose ps -q db) \
 psql -U updoc updoc -c "select flag_id,page_id,status from flags;"
# Additional documentation will be added here
