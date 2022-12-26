# Proof of Work as protection against DDOS attacks

Here is an example of how Proof of Work can be used to protect TCP server agains DDOS attacks by pushing client to make resource-expensive calculations before receiving access to the server functionality.

Hashcash algorithm with sha-2(sha256) hash func was used. It is widely used by email server to protect agains spammers. Also same algorithm used in the most famous cryptocurrency ecosystem - Bitcoin. Bitcoin utilises hashcash algorithm in the new block mining process.

# Whats inside

There are client and server executables in cmd/client and cmd/server respectable. Both wrapped into Docker images and could be easily run via simple Makefile target.

The demo consist of:
* run both server in client in separate docker network
* client will connect to the server.
* server will generate new challenge and send it to the client
* client will solve the challenge and send response back
* server will validate the response for the challenge and if OK send randomly selected quote
* client will read the quote, print it to stdout and exit.
* docker cleanup

# Protection details

Server generates random salt on startup which is added to the hash and validated later, which prevents solution substitution from one server node to another. Also for every new challenge we are adding timestamp to the hash. This helps to protect from the case, when randomly generated nonce was not so random as we expected.

Server sends to the client only hash and a target number, which client should reach by adding the growing counter to the hash. Client doesnt receive neither nonce or timestamp, so it can't store the nonce, target and the calculated solution for later re-use or precalculate them for later.
Also even if nonce is not random - we will not send the same hash for the same nonce, as hash contains hash creation timestamp too and will be different every time.

This trick with nonce encapsulation makes the replay attack almost impossible. 

To consider: We could also add message TTL: 
* add a timestamp to the challenge we send to the client. 
* Then sign entire message and send this signature to the client too.
* Client will have to sand back the whole signed challenge, signature and a solution.

This will add a stable protection from replay attacks, but will increase the network load dramatically as data size we send and receive grows a lot.  

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