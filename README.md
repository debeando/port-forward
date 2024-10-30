# DeBeAndo Zenit Port Forward

Port forward over SSH, allow connect to remote server over SSH to local or private kubernetes cluster.

## Image Description

This image is maintained by DeBeAndo and will be updated regularly on best-effort basis. The image is based on Alpine Linux and only contains the build result of this repository.

## Run

To run container:

```bash
docker run \
	--name zenit-port-forward \
	--publish 3306:3306 \
	--env SSH_HOST="<ssh_host>" \
	--env SSH_KEY="`cat /Users/<username>/.ssh/<private>.pem | base64`" \
	--env REMOTE_HOST="<mysql_host>" \
	debeando/zenit-port-forward
```
