

dft:		# List the targets in this makefile by invoking make
	@echo The default make is to show the targets made by this Makefile and not actually do anything 
	@cat Makefile | grep '^[a-zA-Z0-9_-]*:'

run:	# Build the main.go and run it
	go run main.go


status-postgres:	# Show postgres status but do not change anything.
	postgres.sh

start-postgres:	# Start the postgres DB container.
	postgres.sh start

stop-postgres:	# Stop the postgres DB container.
	postgres.sh stop

migrateup:		# Initialise simple_bank db in postgres container 
	migrate -path  db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose up 
	# docker exec -i postgres createdb --username=root --owner=root simple_bank

migratedown:	# Drop the simple_bank db that was set up above.
	migrate -path  db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose down 
	# docker exec -i postgres dropdb simple_bank

list-schema:	# Display the schema for postgres (db/postgress-accounts.sql) and schema image.
	@cat db/postgres-accounts.sql
	@open images/schema.png

connect:	# Connect to postgres in the docker container.
	@echo Connecting to docker postgres ... use backslash q to quit when done...
	docker exec -it postgres psql

downloads: d1 d2 d3 # The following are the list of installs and downloads.

d1:		# Pull the postgres image from docker.
	./checkit.sh "docker pull postgres:16.2-alpine3.19" "Pull docker image for postgres"

d2:		# Install the psql client (as opposed to TablePlus which costs 90$ or free but very restricted).
	brew install libpg
	@echo "You need to configure your PATH to have the "

d3:		# Install golang migrate to handle database up / down calls - see below.
	brew install golang-migrate

d4:		# Re-initialize the migration scripts used by golang migrate feature in this project.
	rm db/migration/0000*
	migrate create -ext sql -dir db/migration -seq init_schema
	cat db/postgres-accounts.sql > db/migration/000001_init_schema.up.sql
	cat db/postgres-drop.sql     > db/migration/000001_init_schema.down.sql

d5:		# Install the sqlc code generator 
	brew install kyleconroy/sqlc/sqlc

d5b:		# Install the docker version of the sqlc code generator
	docker pull kjconroy/sqlc
	@echo "To run the docker version:"
	@echo "docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate"	


sqlc:		# Generate sqlc CRUD code using sqlc.yaml
	sqlc generate

#q-accounts:	# Init the query/accounts.sql needed the sqlc generation of the accounts CRUD Golang code
#	initc.sh Account Accounts -table accounts -fields 'owner, balance, currency' -values '$1, $2, $3' \
# -update 'balance=$2' # -output db/query/account.sql

# The following 2 bootstraps would be an approximation only .. so they are ignored and not correct
# It shows that the sqlc is pretty good at bootstrapping but attempt to autogenerate the 
#q-entries:	# Init the query/entries.sql needed for sqlc generation of the entries CRUD Golang code 
#	initc.sh Entry Entries -table entries -fields 'account_id, amount' -values '$1, $2' # -update 'account_id=$2, amount=$3' -output db/query/entry.sql

#q-transfers:	# Init the query/transfers.sql needed for sqlc generation of the transfers CRUD in Golang code 
#	initc.sh Transfer Transfers -table transfers -fields 'from_account_id, to_account_id, amount' -values '$$1, $$2, $$3' 
#	  -update 'from_account_id=$$2, to_account_id=$$3, amount=$$4' # -output db/query/transfer.sql

