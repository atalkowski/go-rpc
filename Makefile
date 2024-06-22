

dft:		# List the targets in this makefile by invoking make
	@echo The default make is to show the targets made by this Makefile and not actually do anything 
	@cat Makefile | grep '^[a-zA-Z0-9_-]*:'

run:	# Build the main.go and run it (same as server)
	go run main.go

test:		# Test all unit tests in the project using verbose and coverage mode.
	go test -v -cover ./...

sqlc:		# Generate sqlc CRUD code using sqlc.yaml
	sqlc generate

server:	# Start the main api server (same as run)
	go run main.go

#	showsql -psql sql "select now()"
#	showsql -psql sql "select * from accounts order by id desc limit 10" 
#	showsql -psql sql "select * from entries order by id desc limit 10" 
#	showsql -psql sql "select * from transfers order by id desc limit 10" 

txns:		# Run SQL for the txn update checks 
	showsql -v -psql sql "select a.id, a.owner as donor, a.balance as abal, \
	  b.owner as recipient, b.balance as bbal, t.created_at as txtime,  t.amount \
	 from transfers t \
	inner join accounts a on t.from_account_id = a.id \
	inner join accounts b on t.to_account_id = b.id \
	order by t.created_at desc limit 10"

pgstuff:# Below are the postgres start stoop options
	@echo See the status-pg, start-pg, stop-pg targets and the postgres.sh script.

docker:	# Launch the docker desktop application
	open -a Docker
	@echo Please wait while Docker Desktop loads ...

status-pg:	# Show postgres status but do not change anything.
	postgres.sh

start-pg:	# Start the postgres DB container.
	postgres.sh start

stop-pg:	# Stop the postgres DB container.
	postgres.sh stop

postgres:	# Start the postgres database version postgres:16.2-alpne3.19.
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecret -d postgres:16.2-alpine3.19 

createdb:	# Create the simple_bank database in postgres.
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:		# Drop the simple_bank database from postgres.
	docker exec -it postgres dropdb simple_bank

migrateup:	# Initialise simple_bank db in postgres container 
	migrate -path  db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose up 
	# docker exec -i postgres createdb --username=root --owner=root simple_bank

migratedown:	# Drop the simple_bank db that was set up above.
	migrate -path  db/migration -database "postgresql://root:mysecret@localhost:5432/simple_bank?sslmode=disable" -verbose down 
	# docker exec -i postgres dropdb simple_bank

list-schema:	# Display the schema for postgres (db/postgress-accounts.sql) and schema image.
	@cat db/postgres-accounts.sql
	@open images/schema.png

connect:	# Connect to the postgres if it is running.
	@echo Connecting to docker postgres ... use backslash q to quit when done...
	docker exec -it postgres psql

proto1:		# Generate the gRPC code for the protobuf definitions.
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	       --go-grpc_out=pb  --go-grpc_opt=paths=source_relative  proto/*.proto 

downloads: # The following are the list of installs and downloads.
	@echo These are the d1 d2 d3.... d7 etc


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

d6:		# Install Postgres go library the at lib/pq
	go get github.com/lib/pq

d7:		# Install stretchr/testify
	go get github.com/stretchr/testify
	go get github.com/stretchr/testify/require

d8:		#Install golang gin
	go get -u github.com/gin-gonic/gin

d9:		# Install go viper for handling config

d10:	# Install go protoc and protoc-gen-go and protoc-gen-go-rpc
	brew install protoc-gen-go
	brew install protoc-gen-go-grpc
	
d5b:		# Install the docker version of the sqlc code generator
	docker pull kjconroy/sqlc
	@echo "To run the docker version:"
	@echo "docker run --rm -v $(pwd):/src -w /src kjconroy/sqlc generate"	

api-create:	# Create a test random account using the API
	callapi.sh create

api-delete:	# Call the delete API for ID env value
	callapi.sh delete $(ID)

api-list: # Call the list API for page_size and page_id defaults (or passed ARGS="-page_size X -page_id Y")
	callapi.sh list $(ARGS)

api-get:	# Call the get API for env ID value
	callapi.sh get $(ID)

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

