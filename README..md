
# Golang + Postgres + Kubernetes + gRPC 
This project is from the Udemy on-line course of the title above.

## Setup, build and run
There is a Makefile and some scripts which identify how to setup installs e.g. for the docker postgres service, etc. 



## Database 
### Postgres
Setting up postgres on docker with a postgres-alpine image was the way suggested. 
Installing and connecting are options available in the Makefile (use command make to see).

### dbdiagram.io for schema easy set up
The dbdiagram.io site allows us to define our schema. See Auxillary notes for URL.
The outputs are install sql scripts - see the db diractory - these were exported from the web site which can export 3 DB flavors - we will be using the postgres-accounts.sql export.

### DB Client installation
Go to https://tableplus.com and download the Mac version (or whatever) of Tableplus.

## Auxilliary Notes
1. See https://dbdiagram.io/d/ for the database setup tool which course advises to use.
2. See https://dbdiagram.io/d/Simple-Bank-663bf7a79e85a46d555ba356 for my Simple Bak  
3. See https://hub.docker.com/postgres .. for nots on configuring the postgres container.
