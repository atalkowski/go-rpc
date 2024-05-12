# go-rpc : Golang + Postgres + Kubernetes + gRPC 
This project is from the Udemy on-line course of the title above.

## Setup, build and run
There is a Makefile and some scripts which identify how to setup installs e.g. for the docker postgres service, etc. 
### How to?
Run the following command and/or view contents of the Makefile:
1. make
This should list output options similar to following:
### Listing targets made by this Makefile
1. dft:		# List the targets in this makefile
2. start-postgres:	# Start the postgres DB container
3. stop-postgres:	# Stop the postgres DB container
4. schema:		# Display the schema for postgres (db/postgress-accounts.sql) and schema image.
5. connect:	# Connect to postgres in the docker container
6. downloads: d1 	# The following are tghe list of installs and downloads
601. d1:		# Pull the postgress image from docker
602. d2:		# Install the psql client (as opposed to TablePlus which costs 90$ or free but useless)
603. d3:		# Install golang migrate 
604. d4:		# Migrate data

## Database 
### Postgres
Setting up postgres on docker with a postgres-alpine image was the way suggested. 
Installing and connecting are options available in the Makefile (use command `make` to see).

### dbdiagram.io for schema easy set up
The dbdiagram.io site allows us to define our schema. See the Auxillary notes below for URL.
The outputs of that process are install sql scripts - see the db diractory - these were exported from the web site which can export 3 DB flavors - we will be using the postgres-accounts.sql export.

### DB Client installation
Go to https://tableplus.com and download the Mac version (or whatever) of Tableplus.
However - please note that Tableplus free version is now very restrictive and costs $90 (or more).
Therefore, I've added a showsql command which can run postgres queries from the command line.

## Auxilliary Notes
1. See https://dbdiagram.io/d/ for the database setup tool which course advises to use - this is quite neat.
2. See https://dbdiagram.io/d/Simple-Bank-663bf7a79e85a46d555ba356 for my Simple Bank setup (or NOT).   
3. See https://hub.docker.com/postgres .. for notes on configuring the postgres container.
