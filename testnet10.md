# Testing on testnet10

First, set up [chia-dev-tools](https://github.com/Chia-Network/chia-dev-tools) first. Make sure your venv is active while using `tibet.py`:

```bash
. ./venv/bin/activate
pip install -r requirements.txt
pip install --extra-index-url https://pypi.chia.net/simple/ chia-internal-custody
```

Next, you'll have to switch to testnet10. There's no point in running a full node (db is comparable in size with that on mainnet), so this guide will use [FireAcademy.io](https://fireacademy.io) instead of a full node.

```bash
chia stop all
chia configure --testnet true
chia start wallet
```

Add some peers from [https://alltheblocks.net/testnet10/peers](https://alltheblocks.net/testnet10/peers) (copy the `chia peer -a` commands, but replace "full_node" with "wallet"). Do this until `chia wallet show` says that your wallet is syncing or synced.

Next, get some TXCH (test XCH, the currency of the testnet) from [here](https://xchdev.com/#!faucet.md) or [here](https://testnet10-faucet.chia.net/).

You can also get some TDBX by accepting offers from [Dexie](https://testnet.dexie.space/offers/TDBX/TXCH):
```bash
chia wallet take_offer -m [fee_amount_in_xch] [raw_offer]
```

Before configuring `tibet.py`, you'll need an API key from [FireAcademy.io](https://fireacademy.io). Create an account and get one for free [here](https://dashboard.fireacademy.io/).

Configure:
```bash
rm config.json # delete prev. config
python3 tibet.py config-node --use-preset testnet10 --fireacademyio-api-key [you-api-key]
python3 tibet.py test-node-config
python3 tibet.py set-router --launcher-id d63637fea544958c0e9ce7b6cab2e517b5910980da7fc1a7a734ce0f2e236cd2
python3 tibet.py sync-pairs
```

Time to play! See [TESTING.md](TESTING.md) to get an idea of the possible commands. Do not forget to use the `--fee` switch since the testnet mempool seems full.

*Note on fees*: If your transaction says 'pending', it means that your previous fee was too low. You need to re-generate the offer (cancel it and re-run the command) with a higher fee and then submit it to the blockchain. Keep in mind that the minimum fee bump is 10 million (10,000,000) mojos.

To get back on mainnet:
```bash
chia stop all
chia configure --testnet false
```