

dft:		# List the targets in this makefile
	@echo Listing targets made by this Makefile
	@cat Makefile | grep '^[a-zA-Z0-9_-]*:'

start-postgres:	# Start the postgres DB container
	postgres.sh start

stop-postgres:	# Stop the postgres DB container
	postgres.sh stop

schema:		# Display the schema for postgres (db/postgress-accounts.sql) and schema image.
	@cat db/postgres-accounts.sql
	@open images/schema.png

connect:	# Connect to postgres in the docker container
	@echo Connecting to docker postgres ... use backslash q to quit when done...
	docker exec -it postgres psql

downloads: d1 	# The following are tghe list of installs and downloads

d1:		# Pull the postgress image from docker
	./checkit.sh "docker pull postgres:16.2-alpine3.19" "Pull docker image for postgres"

d2:		# Install the psql client (as opposed to TablePlus which costs 90$ or free but useless)
	brew install libpg
	@echo "You need to configure your PATH to have the "

d3:		# Install golang migrate 
	brew install golang-migrate

d4:		# Migrate data
	migrate create -ext sql -dir db/migration