docker run --name godb -e POSTGRES_USER=golang_db_user -e POSTGRES_PASSWORD=golang_db_password -e POSTGRES_DB=godb -p 7530:5432 -d postgres
docker start godb
docker stop godb
