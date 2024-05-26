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
6. downloads: # The following are tghe list of installs and downloads
601. d1:		# Pull the postgress image from docker
602. d2:		# Install the psql client (as opposed to TablePlus which costs 90$ or free but useless)
603. d3:		# Install golang migrate 
604. d4:		# Migrate data
605. d5:    # Install the sqlc code generator 
606. d6:    # Install the docker version of the sqlc code generator    

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

### Various Go libraries for our CRUD SQL operations
There are 4 contenders for the way to go described here:
1. Database/sql : very fast and straightforward - but need to manually map SQL fields to variables which is error prone.
2. GORM : not fast, CRUD is already implemented => concise code; need to learn queries using gorm's sql syntax.
3. SQLX : quite fast, easy to use, fields mapping via query-text/struct-tags; still need to write the SQL manually.
4. SQLC : very fast, easy to use, auto code-generation from SQL tables defns and relations. Supports POSTGres only atm.
Given our use of POSTGres - we will choose last option SQLC.

### SQLC usage.
Once the sqlc command is install you can set up a sql file input and then autogenerate CRUD code:
- done by passing simple DDL, DML and SQL statements to sqlc.
Example:
  CREATE TABLE authors {
    id BIGSERIAL PRIMARY KEY,
    name text    NOT NULL,
    bio text
  };

  -- name: GetAuthor :one
  SELECT * FROM authors
  WHERE id = $1 LIMIT 1;

  -- name: ListAuthors :many
  SELECT * FROM authors
  ORDER BY name;

  -- name: CreateAuthor :one
  INSERT INT authors ( name, bio ) VALUES ( $1, $2 )
  RETURNING *;

  -- name: DeleteAuthor :exec
  DELETE FROM authors WHERE id = $1;

See more at http://sqlc.dev

## Auxilliary Notes
1. See https://dbdiagram.io/d/ for the database setup tool which course advises to use - this is quite neat.
2. See https://dbdiagram.io/d/Simple-Bank-663bf7a79e85a46d555ba356 for my Simple Bank setup (or NOT).   
3. See https://hub.docker.com/postgres .. for notes on configuring the postgres container.
