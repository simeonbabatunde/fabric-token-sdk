# I Owe yoU Token eXtended (IOUTX)

In this Section, we present a sample that shows how to perform token settlement operations based on the FSC IOU sample extended for three participants (Alice, Bob, Charlie). We utilize a digital USD ($) as our `fungible tokens` for this sample. This example showcases how to mint/create digital USD tokens, and how the minted tokens can be issued or transferred among network participants.

The following business parties are involved in the settlement (token transaction) process:
- `Issuer`: This entity represents the Central Bank, it creates/mints/issues tokens which can be delivered to other parties in the network.
- In this sample, the Issuer mints 10USD and issues everything to `Alice`, who then transfers 7USD to `Bob`. Bob then transfers 2USD to `Charlie`.
- `Auditor`: This entity audits/endorses token transactions.

Each party is running a Smart Fabric Client node with the Token SDK enabled.
The parties are connected in a peer-to-peer network established and maintained by the nodes.


The following are the different token operations that can be performed on a transaction:

## Issuance

Issuance is a business interactive protocol among two parties: an `issuer` of a given token type
and a `recipient` that will become the owner of the freshly created token. [`issue.go`](./views/issue.go) shows a `view` representing the issue operations executed by the Issuer's FSC node.

```go
// IssueCash contains the input information to issue a token
type IssueCash struct {
	// IssuerWallet is the issuer's wallet to use
	IssuerWallet string
	// TokenType is the type of token to issue
	TokenType string
	// Quantity represent the number of units of a certain token type stored in the token
	Quantity uint64
	// Recipient is an identifier of the recipient identity
	Recipient string
}
```

The recipient FSC node executes the `view` in [`accept.go`](./views/accept.go) upon receiving a message from the issuer.  


## Transfer

Transfer is a business interactive protocol among at least two parties: a `sender` and one or more `recipients`. The sender's FSC node executes the `view` representing transfer in [`transfer.go`](./views/transfer.go)  

```go
// Transfer contains the input information for a transfer
type Transfer struct {
	// Wallet is the identifier of the wallet that owns the tokens to transfer
	Wallet string
	// TokenIDs contains a list of token ids to transfer. If empty, tokens are selected on the spot.
	TokenIDs []*token.ID
	// TokenType of tokens to transfer
	TokenType string
	// Quantity to transfer
	Quantity uint64
	// Recipient is the identity of the recipient's FSC node
	Recipient string
	// Retry tells if a retry must happen in case of a failure
	Retry bool
}
```

The recipient FSC node executes the same `view` in [`accept.go`](./views/accept.go) upon receiving a message from the sender.  

Once the transaction is finalized:
- The token spent will disappear from the sender's vault.
- The recipient's vault will contain a reference to the freshly created tokens originated from the transfer (UTXO model). The recipient can query the vault, or the wallet used to derive the recipient identity.


## Queries

Participating nodes (Alice, Bob, Charlie) can execute the `view` in [`list.go`](./views/list.go) to show the list of unspent tokens in a wallet.


Also, the  `view` in [`history.go`](./views/hystory.go) can be executed to show the list of tokens issued by the issuer.


## Testing

To run this sample, we needs to deploy both the `Fabric ` and the `Fabric Smart Client` networks. And then invoke the `views` on the smart client nodes.


### Networks topology (Fabric and FSC)

For Fabric, we will use a simple topology with:
1. Two organizations: Org1 and Org2;
2. Single channel;
2. Org1 runs/endorse the Token Chaincode.

For the FSC network, we have a topology with a node for each business party.
1. `Issuer` and `Auditor` have an Org1 Fabric Identity;
2. `Alice`, `Bob`, and `Charlie` have an Org2 Fabric Identity.

The network topology is described programmatically in [`fabric.go`](./topology/fabric.go) 


### Boostraping the networks

To build the sample, run the following command from the folder: `$GOPATH/src/github.com/hyperledger-labs/fabric-token-sdk/samples/ioutx`.

```shell
go build -o ioutx
```
Run the `ioutx` sample as follows:

``` 
./ioutx network start --path ./testdata
```

The above command will start the Fabric network and the FSC network, and store all configuration files under the `./testdata` directory.
The CLI will also create the folder `./cmd` that contains a go main file for each FSC node. 
These go main files  are synthesized on the fly, and
the CLI compiles and then runs them.

If everything is successful, you will see something like the following:

```shell
2022-02-09 14:17:06.705 UTC [nwo.network] Start -> INFO 032  _____   _   _   ____
2022-02-09 14:17:06.705 UTC [nwo.network] Start -> INFO 033 | ____| | \ | | |  _ \
2022-02-09 14:17:06.705 UTC [nwo.network] Start -> INFO 034 |  _|   |  \| | | | | |
2022-02-09 14:17:06.705 UTC [nwo.network] Start -> INFO 035 | |___  | |\  | | |_| |
2022-02-09 14:17:06.705 UTC [nwo.network] Start -> INFO 036 |_____| |_| \_| |____/
2022-02-09 14:17:06.705 UTC [fsc.integration] Serve -> INFO 037 All GOOD, networks up and running...
2022-02-09 14:17:06.705 UTC [fsc.integration] Serve -> INFO 038 If you want to shut down the networks, press CTRL+C
2022-02-09 14:17:06.705 UTC [fsc.integration] Serve -> INFO 039 Open another terminal to interact with the networks
```

To shut down the networks, just press CTRL-C.

If you want to restart the networks after the shutdown, you can just re-run the above command.
If you don't delete the `./testdata` directory, the network will be started from the previous state.

Before restarting the networks, we can modify the business views to add new functionalities, to fix bugs, and so on.
Upon restarting the networks, the new business views will be available.

To clean up all artifacts, we can run the following command:

```shell
./ioutx network clean --path ./testdata
```

The `./testdata` and `./cmd` folders will be deleted.

## Invoke the business views
Before issuing a 10USD digital token to `Alice`, let us first confirm the wallet balance on both the issuer and recipient.

The following command queries the `issuer`'s wallets for the history of tokens minted/issued so far:

```shell
./ioutx view -c ./testdata/fsc/nodes/issuer/client-config.yaml -f issued -i "{\"TokenType\":\"USD\"}"
```

The next command queries `Alice`'s wallet for the list of unspent tokens:

```shell
./ioutx view -c ./testdata/fsc/nodes/alice/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```
### Token issuance

This command mints 10USD tokens and issues it to `Alice`:

```shell
./ioutx view -c ./testdata/fsc/nodes/issuer/client-config.yaml -f issue -i "{\"TokenType\":\"USD\", \"Quantity\":10, \"Recipient\":\"alice\"}"
```

If everything is successful, you will see something like the following:

```shell
"e8d10cc8bb1155f17b1900d92c2c1874ced832bb5a8ce86f9d4cdc99587d2fe8"
```
Which is the ID of the transaction that issued the fungible token.

We can now query both wallets again for history of issued tokens (`issuer`) and list of unspent tokens (`Alice`).

```shell
./ioutx view -c ./testdata/fsc/nodes/issuer/client-config.yaml -f issued -i "{\"TokenType\":\"USD\"}"
```

```shell
./ioutx view -c ./testdata/fsc/nodes/alice/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```

The list of unspent tokens of type `USD` in `Alice`'s wallet look something like this (beautified): 
```shell
{
  "tokens": [
    {
      "Id": {
        "tx_id": "e8d10cc8bb1155f17b1900d92c2c1874ced832bb5a8ce86f9d4cdc99587d2fe8"
      },
      "Owner": {
        "raw": "MIIEmxMCc2kEggSTCgxJZGVtaXhPcmdNU1ASggkKIAHCT4uQAPEP1883dXxsJYXshm+r/8Sl+KWKQBPtsDMtEiAN1qnd8QKF2YSv6Bftzt1v5XqH30yzJqlzp777FoBRsxpZCgxJZGVtaXhPcmdNU1ASJ2RlZmF1bHQtdGVzdGNoYW5uZWwtZGVmYXVsdC5leGFtcGxlLmNvbRogAvYocoRRzxnU/3a5uT0Lo/0S+q6nBi5/kLf96ay1gywiEAoMSWRlbWl4T3JnTVNQEAEqzgcKRAogDCvVIZKf3BjWSglYHs9hpmIYFoivU2Ny/ikuMz55VmkSICGZg9JMY+E2eGDpXMowMIW/q6qr5f6ruBHewxDeonVhEkQKICqW2eFFYce3rL+W+vu4YEW83SBWpBUae3ZXdtpytC/BEiAO7VqjyDybt5G4yIkRnZt0MaizCkXEIrEcTF5GhEdYdBpECiAAboYv9DUo8rfC2Nn9tF2YoyVh35YrYX60mjfMS/ugGxIgBRL1EIkE2PWm/prA1S4hAwrmHRnnF2FyrKkK1HJVU3siIBB/5FVWOQnly0LQYeW8Xhm1Y1zNFU4YAZ3n4Pn+W1B6KiAlZdXIel/1SQMDjOwUTKJCenIa8gl2Tyg2p3jWtNX/AjIgELICUsvWTp74zgOVGpCdVpYllbLF19jTPup0sDEQHRc6IB5z48Dokk3HGmfT555wnvwLELoid1lmmRlOmlhSfo3rQiANhc+0gpNXAS2XVoZ+fWcROim/cCK1b/f25/YJp9+42EogEsE+x5fHs+4n4+ve8Lnz8UBdvdEYTdfHNQfjETofOr1SIA6Fc43xTun9EYUKuDIpPdVqLli9BilsmHmhhyKv2U6HUiAUYOPpVIbr4PwDpEd5QdET9bu+qIZTx1SQEsq3Ky8TTlogIDYq3UNHAY1tdC0+fiVUM5c2WWBcehqcavMlTDuQxWpiRAogAcJPi5AA8Q/Xzzd1fGwlheyGb6v/xKX4pYpAE+2wMy0SIA3Wqd3xAoXZhK/oF+3O3W/leoffTLMmqXOnvvsWgFGzaiAl/5bWO/zuNCcqEUqLQC7L8ySTQSg5e7w9jP+IZsNc23KIAQogGY6Tk5INSDpyYL+3MftdJfGqSTM1qecSl+SFt67zEsISIBgA3u8SHx52QmoAZl5cRHlnQyLU917a3UbevVzZkvbtGiAJBonQWF/wdeyema1pDDOVvEsxM3CzjvNVrNrc0SKXWyIgEshepduMbetKq3GAjctAj+PR52kMQ9N7TObMAWb6fap6ZzBlAjEA3YJerRXWIkSqyNJj9GWNV1iGDgRaR9F8uSP1ogYdmY0ORg8zIgwntDg41rS7E+JSAjAQowDJ4JD5jNx2O4AFd8JqBQd0BcrEfw+QdlFVhqiHG8wjX67/cT369jPpZtajxdCKAQCSAWgKRAogB3S14S2mIIQOadbFeGXNFOJjxTYcMdMDgO0wxm6M6IISIBp8N5qfXX941p2f4jklZ9lOe07TaD/XvPq9MgljWUikEiAKUJC0MZGIHvJf9lkIHgaZlS00Qi5S30vWvId1rLeqeg=="
      },
      "Type": "USD",
      "Quantity": "0xa"
    }
  ]
}
```

### Token transfer

`Alice` can now transfer `7USD` tokens to `Bob` using:

```shell 
./ioutx view -c ./testdata/fsc/nodes/alice/client-config.yaml -f transfer -i "{\"TokenType\":\"USD\", \"Quantity\":7, \"Recipient\":\"bob\"}"
```

Now, let's check `Alice` and `Bob`'s wallets to see if they are up-to-date.

Alice:

```shell
./ioutx view -c ./testdata/fsc/nodes/alice/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```

Bob:

```shell
./ioutx view -c ./testdata/fsc/nodes/bob/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```

Finally, `Bob` can now transfer `2USD` tokens to `Charlie` using:

```shell 
./ioutx view -c ./testdata/fsc/nodes/bob/client-config.yaml -f transfer -i "{\"TokenType\":\"USD\", \"Quantity\":2, \"Recipient\":\"charlie\"}"
```

Let's query `Bob` and `Charlie`'s wallets to see if they are up-to-date.

Bob:

```shell
./ioutx view -c ./testdata/fsc/nodes/bob/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```

Charlie:

```shell
./ioutx view -c ./testdata/fsc/nodes/charlie/client-config.yaml -f unspent -i "{\"TokenType\":\"USD\"}"
```

## GUI Application
We create a [`gui`](./gui/main.go) application to provide a simplified and more user friendly way to execute the ioutx sample commands against the token-based fabric smart client network. The gui is developed using [`fyne`](https://developer.fyne.io)-a Go based toolkit for building cross platform gui applications.

### Launch GUI APP
From the directory `iouts/gui` run the follwoing command
```shell
go run main.go
```
But first ensure the FSC network is up and runnning before interracting with the gui app.