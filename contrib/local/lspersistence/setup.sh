#!/bin/sh

sudo rm -rf ~/.pstaked/
sudo rm -rf ~/.pstaked/

export BINARY=pstaked
export HOME_APP=$HOME/.pstaked
export CHAIN_ID=localnet

mnemonic="together chief must vocal account off apart dinosaur move canvas spring whisper improve cruise idea earn reflect flash goat illegal mistake blood earn ridge"
mnemonic1="marble allow december print trial know resource cry next segment twice nose because steel omit confirm hair extend shrimp seminar one minor phone deputy"
mnemonic2="axis decline final suggest denial erupt satisfy weekend utility fortune dry glory recall real other evil spatial speed seek rubber struggle wolf tortoise large"
mnemonic3="knock board dolphin cricket strike sense throw security mistake link ocean educate merit pet public economy embark shoot horror pond budget rent toe frozen"


$BINARY init $CHAIN_ID --chain-id $CHAIN_ID
echo $mnemonic | $BINARY keys add val1 --keyring-backend test --recover
echo $mnemonic1 | $BINARY keys add user1 --keyring-backend test --recover
echo $mnemonic2 | $BINARY keys add user2 --keyring-backend test --recover
echo $mnemonic3 | $BINARY keys add user3 --keyring-backend test --recover
$BINARY add-genesis-account $($BINARY keys show val1 --keyring-backend test -a) 10000000000000000000stake
$BINARY add-genesis-account $($BINARY keys show user1 --keyring-backend test -a) 10000000000000000000stake
$BINARY add-genesis-account $($BINARY keys show user2 --keyring-backend test -a) 10000000000000000000stake
$BINARY gentx val1 50000000000000000stake --chain-id $CHAIN_ID --keyring-backend test
$BINARY collect-gentxs


platform='unknown'
unamestr=`uname`
if [ "$unamestr" = 'Linux' ]; then
   platform='linux'
fi

# Enable API and swagger docs and modify parameters for the governance proposal and
# inflation rate from 13% to 33%
if [ $platform = 'linux' ]; then
	sed -i 's/enable = false/enable = true/g' $HOME_APP/config/app.toml
	sed -i 's/swagger = false/swagger = true/g' $HOME_APP/config/app.toml
	sed -i 's%"amount": "10000000"%"amount": "1"%g' $HOME_APP/config/genesis.json
	sed -i 's%"quorum": "0.334000000000000000",%"quorum": "0.000000000000000001",%g' $HOME_APP/config/genesis.json
	sed -i 's%"threshold": "0.500000000000000000",%"threshold": "0.000000000000000001",%g' $HOME_APP/config/genesis.json
	sed -i 's%"voting_period": "172800s"%"voting_period": "30s"%g' $HOME_APP/config/genesis.json
  sed -i 's%"inflation": "0.130000000000000000",%"inflation": "0.330000000000000000",%g' $HOME_APP/config/genesis.json
else
	sed -i '' 's/enable = false/enable = true/g' $HOME_APP/config/app.toml
	sed -i '' 's/swagger = false/swagger = true/g' $HOME_APP/config/app.toml
	sed -i '' 's%"amount": "10000000"%"amount": "1"%g' $HOME_APP/config/genesis.json
	sed -i '' 's%"quorum": "0.334000000000000000",%"quorum": "0.000000000000000001",%g' $HOME_APP/config/genesis.json
	sed -i '' 's%"threshold": "0.500000000000000000",%"threshold": "0.000000000000000001",%g' $HOME_APP/config/genesis.json
	sed -i '' 's%"voting_period": "172800s"%"voting_period": "30s"%g' $HOME_APP/config/genesis.json
  sed -i '' 's%"inflation": "0.130000000000000000",%"inflation": "0.330000000000000000",%g' $HOME_APP/config/genesis.json
fi

# Start
$BINARY start