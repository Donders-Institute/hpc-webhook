#!/bin/bash

set -e

hpc_webhook_key_dir=../configs/

ssh-keygen -t rsa -C "hpc-webhook-server" -f $hpc_webhook_key_dir/hpc-webhook