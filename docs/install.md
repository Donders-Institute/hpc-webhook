# QaaS Installation Instructions

## Obtain the source code

Change to your `GOPATH`, for example on Windows:
```console
$ cd C:\Users\YOURUSERNAME\go\src\github.com\Donders-Institute
```

Obtain the source code:
```console
$ git clone https://github.com/Donders-Institute/hpc-qaas.git
```

Go into the directory:
```console
$ cd hpc-qaas
```

## Configuration

Go to the `configs` folder, 
copy the `qaas-database.env.example` file to `qaas-database.env`, 
and change the contents:

```
# Qaas server settings
QAAS_HOST=qaas.dccn.nl
QAAS_PORT=5111
HOME_DIR=/home
DATA_DIR=/data
KEY_DIR=/keys
INPUT_PRIVATE_KEY_FILE=qaas
INPUT_PUBLIC_KEY_FILE=qaas.pub
PRIVATE_KEY_FILE=qaas
PUBLIC_KEY_FILE=qaas.pub

# Relay computer node settings
RELAY_NODE=relaynode.dccn.nl
RELAY_NODE_TEST_USER=dccnuser
RELAY_NODE_TEST_USER_PASSWORD=somepassword

# Database settings
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=someuser
POSTGRES_PASSWORD=somepassword
POSTGRES_DATABASE=somedatabasename
PGDATA=/data/postgres
```

## Start the services

Run the `start.sh` script in the `scripts` folder.


## Run the tests

Run the `start_test.sh` script in the `test/scripts` folder.
