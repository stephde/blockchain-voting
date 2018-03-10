# Anonymous Voting on Hyperledger Fabric

As part of the course [Building Scalable Blockchain Applications with Big Data Technology](https://hpi.de/naumann/teaching/teaching/ws-1718/building-scalable-blockchain-applications-with-big-data-technology-ps-master.html) at the Hasso-Plattner-Institute, this project aimed at implementing a voting application on the blockchain.

It features a private blockchain using the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric) framework and a Node frontend for administration and voting. Both the used protocols and the frontend are heavily inspired by the [Open Vote Network](https://github.com/stonecoldpat/anonymousvoting), an open-source implementation of anonymous voting on [Ethereum](https://www.ethereum.org/).

## Getting started

1. Run `./setup.sh` to install platform-specific Docker images.
    **macOS:** run `brew tap hyperledger/fabric && brew install fabric-tools`.
2. In `/hyperledger`: Run `./generate` to setup cryptographic material and `./start.sh` to start Docker containers and the Hyperledger network.
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

Unfortunately, we had to disable all cryptographic functionality in it due to severe problems with Hyperledger's way of including libraries in Smart Contracts. Unit tests of the chain code run perfectly locally, but when trying to instantiate it on the blockchain, Hyperledger does not copy the needed library für elliptic curve cryptography into the temporary Docker machine for chain code execution.

*todo: maybe add one or two more details/sources for the issue?*

### Frontend

Our frontend (located in `/votingApp`) **may be used with cautions** as it only serves as a **proof-of-concept**. It is based on [anonymousvoting](https://github.com/stonecoldpat/anonymousvoting)'s frontend but was heavily modified for the purposes of our course's final presentation.

## Benchmarks

### How tun run benchmarks

*todo*

### Results

*todo*





