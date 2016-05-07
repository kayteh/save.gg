# Development Setup

## Required packages

- docker 1.11.1

- git

- go 1.6

- nvm (use node 5.10)

- gulp-cli

## Setting the stage...

```bash
# Clone into the right place
cd $GOPATH/src
mkdir save.gg
cd save.gg
git clone git@github.com:kayteh/save.gg sgg

# Setup a postgres container
docker run -e POSTGRES_PASSWORD=19216801 -e POSTGRES_DB=sgg -e POSTGRES_USER=sgg -e POSTGRES_INITDB_ARGS="-A trust" --name sgg-dev-pg -p 5432:5432 -d postgres

# Setup a NATS container
docker run -d -p 4222:4222 -p 8222:8222 -p 6222:6222 --name sgg-dev-nats nats

# Get govendor
go get -u github.com/kardianos/govendor

# Get dependencies and build
govendor get ./...

cd client
npm i

# Initial migration
sgg-tools migrate

# Make your user
sgg-tools debug-user register -a
```

## Making new things...
```bash
# Start sgg-dev
sgg-dev

# OR

# Start a save.gg cluster
# You only need this if you're working on cluster communications or sgg-router.
sgg-tools dev-cluster

# Start gulp
cd client
gulp

# OR

# Start a gulp with the actual client for dev
cd client
gulp client:dev
```
