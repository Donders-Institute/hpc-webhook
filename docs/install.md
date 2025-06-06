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

copy the `env.template` file to `.env` and build the container with

```bash
$ ./build.sh
```

## Generate the server SSH keys

Run the `generate-keys.sh` script in the `scripts` folder.

## Start the services

Run the `start.sh` script in the `scripts` folder.

## Run the tests

Run the `start_test.sh` script in the `test/scripts` folder.
