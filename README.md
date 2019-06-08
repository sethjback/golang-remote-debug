golang-remote-debug
---

Proof of concept files for remote debugging golang code inside docker container.

Build Either the `Dockerfile.nonbuilt` or `Dockerfile.built` conainer:

```
docker build -t wizfind -f Dockerfile.nonbuilt .
```

Run the container with the correct params:

```
docker run -p 40000:40000 -p 8080:8080 --security-opt=seccomp:unconfined --name wizfind wizfind
```

The `--security-opt=seccomp:unconfined` is necessary to allow dlv to acces the golang proc

Use the `Launch Remote` debug configuration to have vs code attach to the remote (inside container) dlv process.

## limitations

This approach is limited in that the application is not actually running until the debugger attaches.

At this piont, then, it almost seems better to use `telepresence` if you are remote debugging in k8s to mimick having your app running in a pod on the remote cluster.