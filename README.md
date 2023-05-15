# gcounter

Toy GCounter CRDT implementation

## About 

Video Demo: https://www.youtube.com/watch?v=Q9IIfNYPeRw&t=62s

I kept hearing about CRDTs and wanted to learn more about them. I
used ChatGPT to basically be my teacher and I'll jot down the 
rough process I used with ChatGPT screenshots for reproducibility and 
integrity.

## CRDT Background (with ChatGPT/GPT4)

Me: what are crdts
> CRDTs, or Conflict-Free Replicated Data Types, are data structures that allow multiple replicas to be updated independently and concurrently without coordination between them, and then merged together without conflicts.
CRDTs are particularly useful in distributed systems where coordination can be expensive or impractical. For example, they can be used in collaborative editing platforms like Google Docs, where multiple users can make changes to a document at the same time and their changes are combined together.
CRDTs can be classified into two types:
> 1. **Operation-based CRDTs (CmRDTs)**: In this case, the system communicates by transmitting operations. If an operation has been delivered to a replica, it is not delivered again.
2. **State-based CRDTs (CvRDTs)**: In this case, the system communicates by transmitting the full state. Merging is idempotent, commutative, and associative.
Each type of CRDT guarantees that all replicas will eventually converge to the same state, even if updates are made at the same time without coordination, as long as all updates are eventually propagated to all replicas.
It's important to note, though, that while CRDTs ensure eventual consistency, they don't help you if you need strong consistency (where you need to be sure of the state of the system at a particular point in time). Furthermore, they generally require more storage and computation than traditional data structures.

Me: give me the implementation spec of a g counter crdt
> > Sure, let's go through the specification of a G-Counter (Grow-Only Counter) CRDT. This is a simplified representation and might vary depending on your exact needs and the programming language you are using. For this specification, let's assume we have a distributed system with `N` nodes.
> 1. **Data Structure**: Each node `i` maintains a vector (an array) of size `N`, where `N` is the number of nodes in the system. Let's call this vector `V[i]`. The value at index `j` in this vector, `V[i][j]`, represents the count according to node `j`. Initially, all elements of the vector are set to 0.
2. **Increment Operation**: When node `i` increments the counter, it increments the value at index `i` in its vector, i.e., `V[i][i]++`.
3. **Merge Operation**: Nodes occasionally communicate their vectors to each other. When node `i` receives a vector `V'` from another node, it merges `V'` with its own vector `V[i]` by taking the element-wise maximum of the two vectors. In other words, for each `j`, `V[i][j] = max(V[i][j], V'[j])`.
4. **Value Query**: To get the current value of the counter, a node `i` sums up all the elements in its vector `V[i]`. The result is the current value of the counter according to node `i`.

> This is a very basic specification and there are many ways it could be optimized or extended. For example, if you know that updates are always propagated in a certain order, you might be able to use a single counter at each node instead of a vector. Or if you need to support decrement operations, you could use a PN-Counter (Positive-Negative Counter) instead of a G-Counter. But this specification should give you a good starting point for understanding and implementing G-Counters.

Me: when is a gcounter used
> A G-Counter, or Grow-Only Counter, is used in situations where a system needs to maintain a count of events across multiple nodes in a distributed system, and it's acceptable for that count to only increase. G-Counters are particularly useful in systems where nodes can be offline or disconnected for periods of time, where network delays are common, or where you need to minimize coordination between nodes.
Here are some examples of where a G-Counter might be used:
> 1. **Counting Page Views**: In a distributed system where multiple servers are handling web traffic, a G-Counter could be used to count the total number of page views. Each server increments its own counter whenever it serves a page, and these counters are occasionally merged together to get the total count.
2. **IoT Devices**: In an Internet of Things (IoT) scenario, devices could be in areas with spotty or unreliable network connections. A G-Counter could be used to count events like sensor readings or device activations, with each device maintaining its own counter that gets merged when a network connection is available.
3. **Distributed Databases**: In a distributed database, a G-Counter could be used to count events like the number of write operations. This can be useful for monitoring system activity or for implementing features like rate limiting.
4. **Real-Time Analytics**: In systems that provide real-time analytics, G-Counters can be used to count events like user clicks or log entries. Each node in the system can process events independently and update its own counter, which can then be merged with counters from other nodes to get the total count.

> Remember, G-Counters provide eventual consistency, meaning that different nodes might have different counts at any given time, but they will eventually agree on the total count once all updates have been propagated. If you need a counter that supports both increment and decrement operations, or if you need strong consistency (all nodes always agreeing on the count), you would need to use a different data structure.

Me: how do you determine when nodes should communicate with each other in a merge operation
> The timing and frequency of communication between nodes for merge operations in a distributed system, especially in the context of CRDTs like G-Counters, can depend on several factors and can be driven by different strategies. Here are a few common ones:

> 1. **Fixed Time Intervals (Periodic Updates)**: Nodes can communicate their states after a certain fixed amount of time has passed. This is a simple approach but may not be efficient if the state changes frequently or rarely.
2. **On State Change**: Nodes can communicate their states whenever a state change (like an increment operation in the case of a G-Counter) occurs. This can ensure that updates are propagated quickly, but could lead to a lot of communication if updates are frequent.
3. **On Request**: Nodes can communicate their states when another node requests it. This can be efficient in terms of network usage, but it could lead to delays if the requesting node has to wait for the other node to respond.
4. **Hybrid Approaches**: Nodes could use a combination of the above strategies, depending on the system's requirements. For example, they could communicate their states periodically, but also immediately after a certain number of state changes have occurred.
5. **Gossip Protocol**: A popular method used in distributed systems, where nodes randomly share their states with other nodes. Over time, all nodes should reach the same state.

> The right strategy depends on the specific requirements of your system, including the rate of state changes, the cost of communication, the requirement for consistency, and the tolerance for latency. In practice, systems often use a combination of these strategies to balance the need for up-to-date state information with the cost of communication.

### Merge Strategy Used
I used the Gossip Protocol. It's implemented as follows:
* a given worker node randomly decides it wants to merge and sends the merge request to the master with a payload of its data as well as all the other nodes' data
* master picks this up and forwards the merge requst to the rest of 
the nodes with the initiator's data
* the rest of the nodes update their counters with the merge strategy 
described in the spec above

I did not deliberately choose a reason for this strategy; instead, I 
just went with it since it seemed relatively straightforward to implement.

## Possible Improvements...?
* Add another CRDT
* Maybe change this to a operation based CRDT and only propagate the 
state change alone
