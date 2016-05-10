# Development Setup

## Required packages

- docker 1.11.1

- docker-compose

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

# Start the needed docker containers
docker-compose up -d

# Get govendor
go get -u github.com/kardianos/govendor

# Get dependencies and build
govendor sync

# Build binaries (do this after every change~)
go install -v ./... 

# Migrate
sgg-tools migrate
sgg-tools migrate influx
sgg-tools migrate rethink
```
