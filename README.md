# PSTAKE-NATIVE

Welcome to the Pstake's (Cosmos) Liquid Staking Platform repository! This repository contains the source code and resources
for `pstake-native`, a system for liquid-staking tokens on the Persistence blockchain.

## About

`pstake-native` is designed to provide users with a straightforward way to stake tokens to validators in
the IBC and ICA enabled blockchain network. This repository hosts the source code for the platform, allowing developers to
contribute, enhance, and customize the platform's functionalities.

## Getting Started

To get started with `pstake-native`, follow these steps:

**Clone the Repository:** Clone this repository using the following command:
```shell
git clone https://github.com/persistenceOne/pstake-native.git 
```

Install Dependencies: Navigate to the repository directory and install:

```shell
cd pstake-native
make install   
```

## Contributing

### Pull Requests and Changelog

We use an automated system to generate changelog entries from pull requests. When creating a PR:

1. Check the `auto-generate changelog` checkbox in your PR description to include your PR in the changelog.
2. Use conventional commit format in your PR title (e.g., `feat: add feature`, `fix: fix bug`) for automatic categorization.
3. Alternatively, add a custom entry in your PR description using the format `/changelog: Your custom entry here`.

The changelog will be automatically updated when your PR is created or edited.

Feel free to reach out if you have any questions, feedback, or if you'd like to contribute to this project. Let's make
liquid-staking a seamless experience!
