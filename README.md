# View IAC assets

View IAC assests as they are configured in:
- Gitlab (yaml file used by ci/cd pipelines)
- VMWare Cloud Directory

## Configuration

Clone the project and create `.env`:

```shell
cat > .env <<EOF
PICARD_USER=
PICARD_PASSWORD=
GITLAB_TOKEN=
GITLAB_BASEURL=
EOF
```

Create a configuration file to hold sources definitions and populate it with your information:

```shell
cp config.tmpl config.json
```

## Run

Create a binary:

```shell
make
```

Start a webserver on the default port (8080):

```shell
make serve

open http://localhost:8080/
```


View assets as json in terminal:

```shell
make

# gitlab
./dist/iac_linux_amd64 gitlab

# vmware cloud directory
./dist/iac_linux_amd64 vcloud
```