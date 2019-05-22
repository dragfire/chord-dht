# CHORD protocol
A scalable peer-to-peer lookup service for internet applications}

## Details
  
For convenience, Chord nodes are implemented as goroutines. Chord nodes communicate asynchronously with other Chord nodes using JSON messages over **zeroMQ** sockets. The IP address and port number of a node's socket is its access point (address).
   
Chord nodes reveive JSON request messages from the coordinator or other Chord nodes and respond to the sender (or reply-to address specified) directly. We assume the time it takes a node to respond to any message is a random variable (with exponential distribution whose mean is a parameter in your program). The JSON request messages among the Chord nodes and the coordinator are as follows:

    {"do": "join-ring", "sponsoring-node": address } instructing the receipient node to join the Chord ring by contacting the (existing) Chord sponsoring node with the given address.
    {"do": "leave-ring" "mode": "immediate or orderly"} instructing the receipient node to leave the ring immediately (without informing any other nodes) or in an orderly manner (by informing other nodes and transferring its bucket contents to others)
    {"do": "stabilize-ring" }
    {"do": "init-ring-fingers" }
    {"do": "fix-ring-fingers" }
    {"do": "ring-notify", "reply-to": address }
    {"do": "get-ring-fingers", "reply-to": address }
    {"do": "find-ring-successor", "reply-to": address}
    {"do": "find-ring-predecessor", "reply-to": address}
    {"do": "put", "data": { "key" : "a key", "value" : "a value" }, "reply-to": address} instructing the receipient node to store the given (key,value) pair in the appropriate ring node.
    {"do": "get", "data": { "key" : "a key" }, "reply-to": address} instructing the receipient node to retrieve the value associated with the key stored in the ring.
    {"do": "remove", "data": { "key" : "a key" }, "reply-to": address} instructing the receipient node to remove the (key,value) pair from the ring.
    {"do": "list-items", "reply-to": address} instructing the receipient node to respond with a list of the key-value pairs stored at its bucket. 

## References:
~~~
@inproceedings{consistenthashing,
  title = {Chord: A scalable peer-to-peer lookup service for internet applications},
  author = {Ion Stoica, Robert Morris, David Karger, M. Frans Kaashoek and Hari Balakrishnan}, 
  booktitle = {Proceedings of the 2001 conference on Applications, technologies, architectures, and protocols for computer communications},
  series = {SIGCOMM '01},
  year = {2001},
  location = {California, USA},
  pages = {149-160}
}
~~~
