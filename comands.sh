# create new table goose
goose -dir=assets/migrations create tablename sql 

# run goose up migrations
goose -dir=assets/migrations sqlite3 dbname.db up

# run app with hot reloading using air
air

# start .air.toml
air init