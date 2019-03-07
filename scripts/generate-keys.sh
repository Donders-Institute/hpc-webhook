#!/bin/bash

set -e

qaas_key_dir=../configs/

ssh-keygen -t rsa -C "qaas-server" -f $qaas_key_dir/qaas