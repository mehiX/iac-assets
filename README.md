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
# comma separated list of endpoints (no spaces)
VCLOUD_ENDPOINTS=
# comma separated list of tenants (no spaces)
VCLOUD_TENANTS=
GITLAB_TOKEN=
EOF
```

## Run

Create a binary:

```shell
make install
```

Start a webserver on the default port (8080):

```shell
make serve

open http://localhost:8080/
```


View assets as json in terminal:

```shell
make install

# gitlab
./dist/iac gitlab

# vmware cloud directory
./dist/iac vcloud
```