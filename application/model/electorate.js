'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const fs = require('fs');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');
const ccpJSON = fs.readFileSync(ccpPath, 'utf8');
const ccp = JSON.parse(ccpJSON);
// Create a new file system based wallet for managing identities.
const walletPath = path.join(process.cwd(), 'wallet');

class Electorate {
    setElectorate(id, voteID) {
        return new Promise(async(resolve, reject) => {
            try {
                const wallet = new FileSystemWallet(walletPath);

                // Check to see if we've already enrolled the user.
                const userExists = await wallet.exists('user1');
                if(!userExists) {
                    console.log('An identity for the user does not exist in the wallet');
                    reject('An identity for the user does not exist in the wallet')
                }

                // Create a new gateway for connecting to our peer node.
                const gateway = new Gateway();
                await gateway.connect(ccp, { wallet, identity: 'user1', discovery: { enabled: false } });

                // Get the network (channel) our contract is deployed to.
                const network = await gateway.getNetwork('mychannel');

                // Get the contract from the network.
                const contract = network.getContract('hf-voting');
                
                // Submit the specified transaction.
                await contract.submitTransaction(
                    'setElectorate',
                    id.toString(),
                    voteID.toString()
                );
                console.log('Transaction has been submitted');

                // Disconnect from the gateway.
                await gateway.disconnect();
                resolve("유권자가 성공적으로 등록되었습니다.")

            } catch (err) {
                console.log(err);
            }
        });
    }

    getElectorate(id) {
        return new Promise(async(resolve, reject) => {
            try {
                const wallet = new FileSystemWallet(walletPath);

                // Check to see if we've already enrolled the user.
                const userExists = await wallet.exists(id);
                if(!userExists) {
                    console.log('An identity for the user does not exist in the wallet');
                    reject('An identity for the user does not exist in the wallet')
                }

                // Create a new gateway for connecting to our peer node.
                const gateway = new Gateway();
                await gateway.connect(ccp, { wallet, identity: id, discovery: { enabled: false } });

                // Get the network (channel) our contract is deployed to.
                const network = await gateway.getNetwork('mychannel');

                // Get the contract from the network.
                const contract = network.getContract('hf-voting');
                
                // Submit the specified transaction.
                const result = await contract.evaluateTransaction('getElectorate', id.toString());
                console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
                
                // Disconnect from the gateway.
                await gateway.disconnect();

                resolve(result.toString() == "0" ? false : true);
            } catch (err) {
                console.log(err);
            }
        });
    }

    isVoted(id) {
        return new Promise(async(resolve, reject) => {
            try {
                const wallet = new FileSystemWallet(walletPath);

                // Check to see if we've already enrolled the user.
                const userExists = await wallet.exists(id);
                if(!userExists) {
                    console.log('An identity for the user does not exist in the wallet');
                    reject('An identity for the user does not exist in the wallet')
                }

                // Create a new gateway for connecting to our peer node.
                const gateway = new Gateway();
                await gateway.connect(ccp, { wallet, identity: id, discovery: { enabled: false } });

                // Get the network (channel) our contract is deployed to.
                const network = await gateway.getNetwork('mychannel');

                // Get the contract from the network.
                const contract = network.getContract('hf-voting');
                
                // Submit the specified transaction.
                const result = await contract.evaluateTransaction('isVoted', id.toString());
                console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
                
                // Disconnect from the gateway.
                await gateway.disconnect();

                resolve(result.toString() == "0" ? false : true);
            } catch (err) {
                console.log(err);
            }
        });
    }
}

module.exports = new Electorate();