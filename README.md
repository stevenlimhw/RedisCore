# RedisCore: Redis In Several Weekends (using Golang)

## Motivations

This is a hobby project to learn more about Golang and concurrency. I am working on it during weekends and will be documenting the things I learned along the way.

I am learning from the "golang genius" himself on YouTube (Anthony GG) by watching his [Redis project series here.](https://www.youtube.com/watch?v=LMrxfWB6sbQ)

- If you are interested in following his project series, feel free to refer to my source code and documentation for guidance.
- However, the code I write will be alot different than that shown on the video, so please take note!

## Week 1: 6 July - 7 July

### Building a basic TCP server

1. The `Start` function will first check for the current peer connections by spinning up goroutines for the `checkPeerConnections` function.

2. Then, it will call the `acceptPeerConnections` function which will create a new connection to `net.Listener` object and calls `handlePeerConnection` which handles the connection for each peer.

3. The `handlePeerConnection` function will create a new `Peer` and add it into the server's `addPeerCh` channel. Then, the `Peer::readMessages` function will that is called will handle reading messages sent to the peer.

In summary, when the TCP server starts, we need to continuously listen for new peer connections and establish them. For each peer connection established, we then read the messages sent over that connection.
