# Development Setup

## Required packages

- docker 1.12

- docker-compose

- git

- go 1.7

- nvm (use node ^7.2)

- make


## Setting the stage...

```bash
# Clone into the right place
cd $GOPATH/src
mkdir save.gg
cd save.gg
git clone git@github.com:kayteh/save.gg sgg

# Start the needed docker containers
docker-compose up -d

# Get glide
curl -sSL https://glide.sh/get | sh

# Get dependencies and build


# Build binaries (do this after every change~)
make

# Migrate
sgg-tools migrate
sgg-tools migrate influx

# Create first user
sgg-tools debug-user register -a
```
