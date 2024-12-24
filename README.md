# Snow Consensus đemp

## Part 1: How to Run

Follow the steps below to set up and run the software:

1. **Requirements:**  go 1.22

2. **Fast installation:**
   ```bash
   go get
   ```
   ```bash
   cp .env.example .env
   ```
   ```bash
   cd .\cmd\
   go run main.go
   ```

3. **Environment:**  ```.env```

- **`CONTEXT_TIMEOUT`**: The time (in milliseconds) to wait for a request before it times out.

- **`SAMPLE_SIZE`**: The total number of random validators selected for each repeat of the subsampling process.

- **`QUORUM_SIZE`**: The minimum number of validators required to form a "sufficient majority" for decision-making.

- **`DECISION_THRESHOLD`**: The number of consecutive rounds needed to reach consensus.

- **`COUNT_NODE`**: The total number of nodes in the network.

4. **Features:**

- Simulates a P2P network where each node listens to a separate port. Nodes communicate with each other via HTTP.

- When a transaction is created, the Snow Consensus algorithm is executed between the nodes to reach a decision.

- A simple demo to validate messages using public keys.
   