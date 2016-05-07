# Development Setup

## Setting the stage...

```bash
# Clone into the right place
cd $GOPATH/src
mkdir save.gg
cd save.gg
git clone git@github.com:kayteh/save.gg sgg

# Setup a postgres container
docker run -e POSTGRES_PASSWORD=19216801 -e POSTGRES_DB=sgg -e POSTGRES_USER=sgg -e POSTGRES_INITDB_ARGS="-A trust" --name sgg-dev-pg -p 5432:5432 -d postgres

# Get govendor
go get -u github.com/kardianos/govendor

# Get dependencies and build
govendor get ./...

# Initial migration
sgg-tools migrate

# Make your user
sgg-tools debug-user register -a
```

## Making new things...
```bash
# Start sgg-dev
sgg-dev

# Start gulp
cd client
gulp

# OR

# Start a gulp with the actual client for dev
cd client
gulp client:dev
```
