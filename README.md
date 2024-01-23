## Design and implement “Word of Wisdom” tcp server.

• TCP server should be protected from DDOS attacks with the Proof of Work (https://en.wikipedia.org/wiki/Proof_of_work),
the challenge-response protocol should be used.

• The choice of the POW algorithm should be explained.

• After Proof Of Work verification, server should send one of the quotes from “word of wisdom” book 
or any other collection of the quotes.

• Docker file should be provided both for the server and for the client that solves the POW challenge

# The choice of POW algorithm

Proof of work (PoW) serves as a cryptographic validation method, wherein one party, known as the prover,
provides evidence to others, the verifiers, that a specific amount of computational effort has been invested. 
Verifiers can easily confirm this effort with minimal exertion on their part.

Here is a list of known proof-of-work functions:

* Merkle tree–based
* Guided tour puzzle protocol
* Hashcash

I decided to use Hashash because I found a lot of documentation on this algorithm,
and it seemed quite simple but effective to me.
The other two algorithms, Merkle tree and Guided tour puzzle, presented certain drawbacks:

* Merkle tree required the server to undertake substantial work in validating the client's solution,
involving multiple hash calculations for each level of the tree.
* Guided tour puzzle necessitated frequent client requests to the server for additional parts of the guide,
complicating the protocol's logic.

In contrast, Hashash is easy to implement, does not require a lot of work to check on the server side,
and it is possible to dynamically change the complexity depending on the load on the server.

The main downside of Hashash is that the difficulty of the puzzle only depends on the
length of the required all-zeros prefix. Adding a single zero to the prefix doubles the number
of attempts that the prover must afford, and also the variance increases exponentially.
# How to

Build docker image for client and server:

```shell
make docker
```

Run example which start client and server. 
Client connects to the server, solves challenge and prints some random quote:

```shell
make example
```