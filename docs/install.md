# HPC-webhook Installation Instructions

## Obtain the source code

Change to your `GOPATH`, for example on Windows:
```console
$ cd C:\Users\YOURUSERNAME\go\src\github.com\Donders-Institute
```

Obtain the source code:
```console
$ git clone https://github.com/Donders-Institute/hpc-webhook.git
```

Go into the directory:
```console
$ cd hpc-webhook
```

## Configuration

Go to the `configs` folder, 
copy the `hpc-webhook-database.env.example` file to `hpc-webhook-database.env`, 
and change the contents:

```
# HPC webhook server settings
HPC_WEBHOOK_HOST=hpc-webhook.dccn.nl
HPC_WEBHOOK_INTERNAL_PORT=5111
HPC_WEBHOOK_EXTERNAL_PORT=443
HOME_DIR=/home
DATA_DIR=/data
PRIVATE_KEY_FILE=/run/secrets/hpc_webhook_private_key
PUBLIC_KEY_FILE=/run/secrets/hpc_webhook_public_key

# Relay computer node settings
RELAY_NODE=relaynode.dccn.nl
CONNECTION_TIMEOUT_SECONDS=30

# Database settings
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=someuser
POSTGRES_PASSWORD=somepassword
POSTGRES_DATABASE=somedatabasename
```

## Generate the server SSH keys

Run the `generate-keys.sh` script in the `scripts` folder.

## Start the services

Run the `start.sh` script in the `scripts` folder.

## Run the tests

Run the `start_test.sh` script in the `test/scripts` folder.
