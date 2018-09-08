# malice-mcafee

[![Circle CI](https://circleci.com/gh/malice-plugins/mcafee.png?style=shield)](https://circleci.com/gh/malice-plugins/mcafee) [![License](http://img.shields.io/:license-mit-blue.svg)](http://doge.mit-license.org) [![Docker Stars](https://img.shields.io/docker/stars/malice/mcafee.svg)](https://store.docker.com/community/images/malice/mcafee) [![Docker Pulls](https://img.shields.io/docker/pulls/malice/mcafee.svg)](https://store.docker.com/community/images/malice/mcafee) [![Docker Image](https://img.shields.io/badge/docker%20image-433MB-blue.svg)](https://store.docker.com/community/images/malice/mcafee)

Malice McAfee AntiVirus Plugin

> This repository contains a **Dockerfile** of [mcafee](https://www.mcafee.com/enterprise/en-us/products/virusscan-enterprise-for-linux.html).

---

### Dependencies

- [ubuntu:bionic (_84.1 MB_\)](https://hub.docker.com/_/ubuntu/)

## Installation

1. Install [Docker](https://www.docker.com/).
2. Download [trusted build](https://store.docker.com/community/images/malice/mcafee) from public [docker store](https://store.docker.com): `docker pull malice/mcafee`

## Usage

```
docker run --rm malice/mcafee EICAR
```

### Or link your own malware folder:

```bash
$ docker run --rm -v /path/to/malware:/malware:ro malice/mcafee FILE

Usage: mcafee [OPTIONS] COMMAND [arg...]

Malice McAfee AntiVirus Plugin

Version: v0.1.0, BuildTime: 20180903

Author:
  blacktop - <https://github.com/blacktop>

Options:
  --verbose, -V          verbose output
  --elasticsearch value  elasticsearch url for Malice to store results [$MALICE_ELASTICSEARCH_URL]
  --table, -t            output as Markdown table
  --callback, -c         POST results back to Malice webhook [$MALICE_ENDPOINT]
  --proxy, -x            proxy settings for Malice webhook endpoint [$MALICE_PROXY]
  --timeout value        malice plugin timeout (in seconds) (default: 120) [$MALICE_TIMEOUT]
  --help, -h             show help
  --version, -v          print the version

Commands:
  update  Update virus definitions
  web     Create a McAfee scan web service
  help    Shows a list of commands or help for one command

Run 'mcafee COMMAND --help' for more information on a command.
```

## Sample Output

### [JSON](https://github.com/malice-plugins/mcafee/blob/master/docs/results.json)

```json
{
  "mcafee": {
    "infected": true,
    "result": "EICAR test file",
    "engine": "5600.1067",
    "database": "9005",
    "updated": "20180903"
  }
}
```

### [Markdown](https://github.com/malice-plugins/mcafee/blob/master/docs/SAMPLE.md)

---

#### McAfee

| Infected |     Result      |  Engine   | Updated  |
| :------: | :-------------: | :-------: | :------: |
|   true   | EICAR test file | 5600.1067 | 20180903 |

---

## Documentation

- [To write results to ElasticSearch](https://github.com/malice-plugins/mcafee/blob/master/docs/elasticsearch.md)
- [To create a McAfee scan micro-service](https://github.com/malice-plugins/mcafee/blob/master/docs/web.md)
- [To post results to a webhook](https://github.com/malice-plugins/mcafee/blob/master/docs/callback.md)
- [To update the AV definitions](https://github.com/malice-plugins/mcafee/blob/master/docs/update.md)

## Issues

Find a bug? Want more features? Find something missing in the documentation? Let me know! Please don't hesitate to [file an issue](https://github.com/malice-plugins/mcafee/issues/new).

## TODO

- [ ] add licence expiration detection

## CHANGELOG

See [`CHANGELOG.md`](https://github.com/malice-plugins/mcafee/blob/master/CHANGELOG.md)

## Contributing

[See all contributors on GitHub](https://github.com/malice-plugins/mcafee/graphs/contributors).

Please update the [CHANGELOG.md](https://github.com/malice-plugins/mcafee/blob/master/CHANGELOG.md) and submit a [Pull Request on GitHub](https://help.github.com/articles/using-pull-requests/).

## License

MIT Copyright (c) 2017 **blacktop**
