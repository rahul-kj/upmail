# Checkup + Email

Provides email notifications for [sourcegraph/checkup](https://github.com/sourcegraph/checkup).

## Usage

```console
$ upmail --help
usage: upmail --recipient=RECIPIENT --smtp-server=SMTP-SERVER --sender=SENDER [<flags>]

Flags:
  --help                         Show context-sensitive help (also try --help-long and --help-man).
  --config="checkup.json"        checkup.json config file location
  --recipient=RECIPIENT          recipient for email notifications
  --interval="10m"               check interval (ex. 5ms, 10s, 1m, 3h)
  --smtp-server=SMTP-SERVER      SMTP server for email notifications
  --sender=SENDER                SMTP default sender email address for email notifications
  --smtp-username=SMTP-USERNAME  SMTP server username
  --smtp-password=SMTP-PASSWORD  SMTP server password
  --debug                        run in debug mode
  --version                      Show application version.
```

## Deploy this on cloudfoundry
* Download the `upmail_linux_amd64` binary from the [releases](https://github.com/rahul-kj/upmail/releases)

* Create a manifest file using the following skeleton:

```
---
applications:
- name: checkup
  path: ./
  instances: 1
  memory: 512M
  command: ./upmail_linux_amd64
  buildpack: binary_buildpack
  no-route: true
  health-check-type: none
  env:
    RECIPIENT_EMAIL:
    SENDER_EMAIL:
    INTERVAL:
    SMTP_SERVER:
    SMTP_USERNAME:
    SMTP_PASSWORD:
```

**NOTE: `SMTP_USERNAME` and `SMTP_PASSWORD` are optional.**

* Create the checkup.json file in the same folder as `upmail_linux_amd64`

```
{
"checkers": [{
        "type": "tcp",
        "endpoint_name": "PWS Website",
        "endpoint_url": "run.pivotal.io:443",
        "attempts": 5,
        "tls": true,
        "tls_skip_verify": true
    },
    {
        "type": "http",
        "endpoint_name": "Pivotal Website",
        "endpoint_url": "https://pivotal.io",
        "attempts": 5,
        "tls_skip_verify": true
    }]
}
```

* Verify your directory contents are as follows:

```console
$ ls
checkup.json		manifest.yml		upmail_linux_amd64
```

* Push the application into cloudfoundry
`cf push`
