# To update the AV run the following:

```bash
$ docker run --name=mcafee malice/mcafee update
```

## Then to use the updated AVG container:

```bash
$ docker commit mcafee malice/mcafee:updated
$ docker rm mcafee # clean up updated container
$ docker run --rm malice/mcafee:updated EICAR
```
