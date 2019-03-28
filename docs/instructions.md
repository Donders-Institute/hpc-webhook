# webhook testing instructions

We propose to use the following steps to try out webhook
([https://github.com/Donders-Institute/hpc-webhook](https://github.com/Donders-Institute/hpc-webhook)).

## 1. Create and debug your bash script on mentat

Login to a mentat machine of choice, for example `mentat005.dccn.nl`.

Create your bash script to be run on the cluster, for example `test.sh`:
```
#!/bin/bash
#PBS -l walltime=00:01:30
#PBS -l mem=10Mb
echo "script start..."
hostname
whoami
env
sleep 60
echo "script stop..."
```
The asked resources should be in the script defined with `#PBS ...` comments.

Place this script in a folder of choice.

Try out your script, by submitting it:
```
$ qsub test.sh
```

Check the status of your cluster job with:
```
$ qstat
```
Make sure your script ran succesfully.
You should end up with two text files, one stdout (`.o`) and one stderr (`.e`).

## 2. Register a webhook on mentat005

Login to `mentat005.dccn.nl` where the new `hpcutil` tools are installed.

Create a new webhook using the the name of your script `test.sh`:
```
$ hpcutil webhook create test.sh 
```

You should get a message like
```
INFO[0000] webhook created successfully with URL: https://hpc-webhook.dccn.nl:443/webhook/5126d168-e3f1-4c7f-b228-a57fbaf007c4
```
Copy this webhook payload URL, we need it later.

## 3. Configuring the webhook client on github.com

Go to your owned Github repository, for example `https://github.com/rutgervandeelen/simple`.

Choose `Settings`.

Choose `Webhooks`.

Fill in the webhook payload URL, for example:
```
https://hpc-webhook.dccn.nl:443/webhook/5126d168-e3f1-4c7f-b228-a57fbaf007c4
```

Update it.

## 4. Commit your software changes to github

Change your software and commit these changes to your github repository.

For example, change the README file and commit it.

Check `Settings > Webhooks` and see if the delivery of the payload was succesful.
There should be a green tickmark.

## 5. Check the results on mentat

Login to a mentat machine of choice, for example `mentat005.dccn.nl`.

Go to your webhook results folder `/home/dccngroup/dccnuser/.webhooks/5126d168-e3f1-4c7f-b228-a57fbaf007c4`.

Run `qstat` to check if your submitted job is queued, running or completed.

Once the script is finished you should have two text files in this result folder, for example:

```
$ cd ~/.webhooks/5126d168-e3f1-4c7f-b228-a57fbaf007c4
$ ls -1

payload
script
test.sh.e34986226
test.sh.o34986226
```
