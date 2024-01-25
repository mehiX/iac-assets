# View IAC assets

View IAC assests as they are configured in:
- Gitlab (yaml file used by ci/cd pipelines)
- VMWare Cloud Directory

## Configuration

Create a configuration file to hold sources definitions and populate it with your information:

```shell
cp config.tmpl config.json
```

For Gitlab get a token [here](https://git.lpc.logius.nl/-/profile/personal_access_tokens) with scope `read_api`.

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