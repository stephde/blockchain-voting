# Anonymous Voting on Hyperledger Fabric

As part of the course [Building Scalable Blockchain Applications with Big Data Technology](https://hpi.de/naumann/teaching/teaching/ws-1718/building-scalable-blockchain-applications-with-big-data-technology-ps-master.html) at the Hasso Plattner Institute, this project aimed at implementing a voting application on the blockchain.

It features a private blockchain using the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) framework and a Node frontend for administration and voting. Both, the used protocols and the frontend, are heavily inspired by the [Open Vote Network](https://github.com/stonecoldpat/anonymousvoting), an open-source implementation of anonymous voting on [Ethereum](https://www.ethereum.org/).

## Getting started

1. Run `./setup.sh` to install platform-specific Docker images.
    **macOS:** run `brew tap hyperledger/fabric && brew install fabric-tools`.
2. In `/hyperledger`: Run `./generate.sh` to setup cryptographic material and `./start.sh` to start Docker containers and the Hyperledger network.
3. In `/votingApp`: Run  `npm install` and `npm start`
4. Open `localhost:3000` in your browser to access the administration and voting interface.

## Background

In this project we set ourselves the goal to build a decentralized, anonymous voting application. While looking for related work, we stumbled upon Patrick McCorry's work on implementing [anonymous voting on the Ethereum blockchain](https://github.com/stonecoldpat/anonymousvoting) without the need for a tally authority.

In the [respective paper](http://fc17.ifca.ai/preproceedings/paper_80.pdf) McCorry lays the cryptographic groundwork so that votes can be both stored publicly *and* anonymously, so that no single vote but only the final tally can be computed. Furthermore, he implemented the so-called Open Vote protocol in Solidity, i.e. Ethereum's programming language for Smart Contracts.

> The Open Vote Network (OV-net) is a 2-round decentralized voting protocol with the following attractive features
>
> - All communication is public - no secret channels between voters are required.
> - The system is self-tallying - no tallying authorities are required.
> - The voter's privacy protection is maximum - only a full collusion that involves all other voters in the election can uncover the voter's secret vote.
> - The system is dispute-free - everybody can check whether all voters act according to the protocol, hence ensuring the the result is publicly verifiable.
>
> – https://github.com/stonecoldpat/anonymousvoting

Open issues of his implementation on Ethereum are:

- For the tally to be computable, all registered voters must have cast their vote.
- Only yes/no-votes are possible.
- Due to the current implementation and Ethereum's blocksize, no more than 50 voters are possible.

Spending most of our time on **porting anonymousvoting's code to Hyperledger Fabric** (i.e. from Solidity to Go), we primarily aimed at running the Open Node Network as is on Hyperledger. This is already benefitial for certain use cases, because a private blockchain can be configured more easily (e.g. increase blocksize to support more users) and is cheaper than Ethereum's running expenses.

## Issues

### Chaincode

The blockchain application uses the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) framework and implements chaincode resembling the Open Vote Network's smart contract, i.e. it exposes the same interface/protocol.

Unfortunately, we had to disable all cryptographic functionality in it due to severe problems with Hyperledger's way of including libraries in Smart Contracts. Unit tests of the chain code run perfectly locally, but when trying to instantiate it on the blockchain, Hyperledger does not copy the needed library for elliptic curve cryptography into the temporary Docker Container that is used for chaincode execution.

In more details the problem is the following:

The elliptic curve library uses C header files that need to be compiled before being used in Go.
_Govendor_ is Hyperledger's package manager and responsible for compiling and loading these dependencies.
In our local development setup this works as expected, however, in Hyperledger these files are not built correctly.
There is a [bug report](https://github.com/ethereum/go-ethereum/issues/2738) for _another_ Go package manager that describes the problem accurately (except for the fact that Hyperledger uses a different package manager).
We assume though that the problems are related.

Due to that the Master branch of this repository only contains chaincode that is _not_ relying on any cryptographic computations.
If you want to have a look at the cryptographic implementation, refer to the `cryptoChaincode` branch.
Be aware that this branch will not work with Hyperledger!

### Frontend

Our frontend (located in `/votingApp`) **may be used with cautions** as it only serves as a **proof-of-concept**. It is based on [anonymousvoting](https://github.com/stonecoldpat/anonymousvoting)'s frontend but was heavily modified for the purposes of our course's final presentation.

## Benchmarks

As part of the final presentation for the course we conducted some performance measurements.
At this point we do not want to provide a in-depth benchmark for Hyperledger, but rather evaluate whether the proposed system is able to handle our use-case.

The first benchmark compares the chaincode implementations with and without cryptography.
Since we were not able to run cryto chaincode on Hyperledger, we did this in unit tests using a mocked Hyperledger instance.
Obviously, a real Hyperledger system will behave differently, so that the absolute execution times are not that meaningful.
They are good enough though to estimate the overhead for doing cryptography.

The second benchmark runs a voting on a Hyperledger instance (without cryptography).
It is a small bash script that uses Hyperledger's CLI.

### How tun run benchmarks

Enable the end-to-end unit tests in the chaincode directory on the respective branch to run a timed vote with and without crypto.

Run `benchmark-voting.sh` to conduct a vote on a Hyperledger instance.
As part of our project we tested this file the Docker containers that Hyperledger provides.

### Results

Since Hyperledger does not provide an implementation of a consensus protocol as of now, we only run the tests with a single Orderer node.
In the dockerized setting we were able to reach a throughput of about 80 transactions per second on a MacBook Pro.
Obviously Hyperledger offers lots of options to tune its performance, so there is a good chance that we might be able to increase the throughput in the future.

We noticed that the throughput was mainly limited by Hyperledger's event hub, which lost some transactions when sending too many transactions in a short amount of time.
Within the span of the project we were not able to fix this properly.

## Next Steps

In a student's project there is usually not enough time to implement _everything_ that is needed to run it in production.
At the end of this project we would like to point out three topics that interested readers may feel free to improve on.

__Implement Consensus__

As of now, Hyperledger does not provide a default implementation of a consensus protocol.
Obviously, the consensus protocol is one of the most critical parts of any blockchain application, so providing an implementation that may be used within different applications would be a big achievement.

With that in mind the next point becomes feasible.

__Create test setup with multiple orderer nodes__

Within our project only scaled the number of peer nodes, but not the number of orderer nodes in our Hyperledger network.
But, having only one orderer node again implies that one has to trust this node.
As soon as there is a consensus implementation, scaling the number of orderer nodes is the logical next step.

__Improve Crypto Implementation in Go__

All the cryptographic implementation that is being used for encrypting the votes and conducting the vote, is taken from an implementation for Ethereum written in Solidity.
We set ourselves the goal to port this code to Go and improve its readability and maintainability.
In sum we did not achieve the latter two ones.
Without an in-depth understanding of what is happening in the cryptographic part it is not possible to alter the source code.
Additionally the lack of understanding makes finding bugs hard.

We encourage interested readers to improve the crypto chaincode, fix existing bugs, and provide a more readable implementation.
