# Proof of Work as protection against DDOS attacks

Here is an example of how Proof of Work can be used to protect TCP server agains DDOS attacks by pushing client to make resource-expensive calculations before receiving access to the server functionality.

Hashcash algorithm with sha-2(sha256) hash func was used. It is widely used by email server to protect agains spammers. Also same algorithm used in the most famous cryptocurrency ecosystem - Bitcoin. Bitcoin utilises hashcash algorithm in the new block mining process.

# How to run

There is list of handful make targets defined in the Makefile. 

The simpliest way to run both the server and the client is with help of docker:
```shell
> make demo
```

This Makefile target will:
* create docker network
* build both server and client images. Single unified Dockerfile used for both binaries.
* start both images. client will print to stdout how it pass the Proof of work challenge and the payload it received from the server. client will exit when finish.
* stop server container
* remove build docker images
* remove docker network

## In consideration
* stdlib logger used for simplicity. in real life i would use zap or one with similar performance structured logger. 
* binary serializations support data not more than 255(byte) length for simplicity.
* There is single Dockerfile that unified and used both for client and server.
* Tests are not covering entire codebase. Just a few examples added to show how i would test some parts.
* There is still room for performance scaling of the TCP server itself. we can use directly linux epoll for example to "offload" the connection where there is nothing to read at the moment.
* There is a lot of room for DDOS protection optimizations. Like algo type rotations per new incoming connection
* There no metrics where added. I would add prometheus library(which adds go metrics by default like memory utilization, etc) and cover incoming connections, their duration, closing status(ok/error), ratio of succeed POW, etc.