golang-remote-debug
---

Proof of concept files for remote debugging golang code inside docker container.

**Disclaimer** - All of this code/process has been tested on Arch Linux using vs-code and a k8s cluster deployed via `kops` running on AWS. If anything is different in your setup ymmv.

## General Process / Theory

The goal is to allow us to run code as if it was in the k8s cluster, attaching our own debugger to it.

There are two things to note right off:

1. The code/process against which you will be doing the debugging is running locally - i.e. you are not attaching a debugger to the code running in the cluster, rather you are runnig debugging code locally and routing cluster traffic transparently to it.
2. `telepresence` handles swapping/routing the k8s traffic to your local code. Debugging of go code proceeds normally (usually `dlv` sitting in the middle of your applicaiton code and the IDE)

To accomplish this we will configure vs-code to "remote" debug our code, namely inside of a docker container. Once this is configured and working, we can use `telepresence` to pass all traffic from the k8s cluster to the container running locally that we are going to debug.

## Building the container 

Build Either the `Dockerfile.nonbuilt` or `Dockerfile.built` conainer:

```
docker build -t wizfind -f Dockerfile.nonbuilt .
```

Run the container with the correct params:

```
docker run -p 40000:40000 -p 8080:8080 --security-opt=seccomp:unconfined --name wizfind wizfind
```

The `--security-opt=seccomp:unconfined` is necessary to allow dlv to acces the golang proc

## Remote Debugging in VS Code

Use the `Launch Remote` debug configuration to have vs code attach to the remote (inside container) dlv process.

Normally vs code will launch it's own `dlv` process and connect to it. This configuration will tell it to connect to the process running on docker

## Telepresence

There is a test deployment/service config in the `k8s` directory under wizFind. Make sure you adjust the namespace, etc. to match your cluster.

Apply using:

```
kubectl apply -f deploy.yaml
```

This will create both a deployment and loadbalancer service to expose it to the internet. To find the external endpoint you can:

```
kubectl get services
```

On AWS the external service will look something like `a925d1da98b8a11e9b51a06dc9dfcc62-1023382574.us-west-2.elb.amazonaws.com`

You can test this by visiting the external IP on port `8080` (the default)

To do local debugging, use telepresence to swap out the in-cluster service with the local docker container.

```
telepresence --swap-deployment wizfind --expose 8080 --docker-run --rm -p 40000:40000 -p 8080:8080 --security-opt=seccomp:unconfined --name wizfind wizfind
```

In vs-code, launch the debugger and you will see the wizFind service start. Going to the external IP noted above, the traffic will now be routed to your locally running container. You can test this by setting breakpoints in vs code and confirming they are triggered when you hit the external URL.

Once you stop debugging in vs code the container will exit and telepresence will restore the original config.

