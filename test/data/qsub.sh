#!/bin/bash
#PBS -l walltime=00:01:30
#PBS -l mem=10Mb
echo "script start..."
hostname
whoami
env
sleep 60
echo "script stop..."
