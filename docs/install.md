# QaaS Installation Instructions

## Obtain the source code

Change to your `GOPATH`, for example on Windows:
```
cd C:\Users\YOURUSERNAME\go\src\github.com\Donders-Institute
```

Obtain the source code:
```
git clone https://github.com/Donders-Institute/hpc-qaas.git
```

Go into the directory:
```
cd hpc-qaas
```

## Configuration

Go to the `configs` folder, 
copy the `qaas-database.env.example` file to `qaas-database.env`, 
and change the contents:

```
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