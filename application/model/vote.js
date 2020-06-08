'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', '..', 'first-network', 'connection-org1.json');

const crypto = require('crypto');

class Vote {
    async setVote(name) {
        try {
            const id = await crypto.randomBytes(32).toString('hex');
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet = new FileSystemWallet(walletPath);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists('user1');
            if(!userExists) {
                console.log('An identity for the user does not exist in the wallet');
                console.log(`Run the registerUser.js application before retrying`);
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork('mychannel');

            // Get the contract from the network.
            const contract = network.getContract('hf-voting');
            
            // Submit the specified transaction.
            await contract.submitTransaction(
                'setVote',
                id.toString(),
                name.toString()
            );
            console.log('Transaction has been submitted');

            // Disconnect from the gateway.
            await gateway.disconnect();
            return "투표가 성공적으로 등록되었습니다.";

        } catch (err) {
            console.log(err);
            return "투표가 정상적으로 등록되지 않았습니다.";
        }
    }

    async getVote(id) {
        try {
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet = new FileSystemWallet(walletPath);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists('user1');
            if(!userExists) {
                console.log('An identity for the user does not exist in the wallet');
                console.log(`Run the registerUser.js application before retrying`);
            }

            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });
            
            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork('mychannel');

            // Get the contract from the network.
            const contract = network.getContract('hf-voting');
            
            // Submit the specified transaction.
            const result = await contract.evaluateTransaction('getVote', id.toString());
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

            // Disconnect from the gateway.
            await gateway.disconnect();
            return result.toString();
        } catch(err) {
            console.log(err)
            return "투표가 정상적으로 조회되지 않았습니다.";
        }
    }

    async getAllVotes() {
        try {
            console.log(ccpPath);
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            console.log(`walletPath: ${walletPath}`);
            const wallet = new FileSystemWallet(walletPath);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists('user1');
            if(!userExists) {
                console.log('An identity for the user does not exist in the wallet');
                console.log(`Run the registerUser.js application before retrying`);
            }
            
            // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork('mychannel');

            // Get the contract from the network.
            const contract = network.getContract('hf-voting');
            
            // Submit the specified transaction.
            const result = await contract.evaluateTransaction('getAllVotes');
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);

            return result.toString();
        } catch(err) {
            console.error(`Failed to evaluate transaction: ${err}`);
            return "투표가 정상적으로 조회되지 않았습니다.";
        }
    }

    async vote(id, voteID, candidateName) {  // Electorate ID, Vote ID, Candidate Name
        try {
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet = new FileSystemWallet(walletPath);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists(id);
            if(!userExists) {
                console.log('An identity for the user does not exist in the wallet');
                console.log(`Run the registerUser.js application before retrying`);
            }

                // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork('mychannel');

            // Get the contract from the network.
            const contract = network.getContract('hf-voting');

            // Submit the specified transaction.
            await contract.submitTransaction(
                'vote',
                id.toString(),
                voteID.toString(),
                candidateName.toString()
            );
            console.log('Transaction has been submitted');
            return '투표하였습니다.';
        } catch (err) {
            console.log(err);
            return '투표 실패하였습니다.';
        }
    }

    async getHistoryByVoteID(voteID) {
        try {
            // Create a new file system based wallet for managing identities.
            const walletPath = path.join(process.cwd(), 'wallet');
            const wallet = new FileSystemWallet(walletPath);

            // Check to see if we've already enrolled the user.
            const userExists = await wallet.exists('user1');
            if(!userExists) {
                console.log('An identity for the user does not exist in the wallet');
                console.log(`Run the registerUser.js application before retrying`);
            }

                // Create a new gateway for connecting to our peer node.
            const gateway = new Gateway();
            await gateway.connect(ccpPath, { wallet, identity: 'user1', discovery: { enabled: true, asLocalhost: true } });

            // Get the network (channel) our contract is deployed to.
            const network = await gateway.getNetwork('mychannel');

            // Get the contract from the network.
            const contract = network.getContract('hf-voting');

            // Submit the specified transaction.
            await contract.evaluateTransaction(
                'getHistoryByVoteID',
                voteID.toString()
            );
            console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
            return result.toString();
        } catch (err) {
            console.log(err);
            return '투표 내역 조회에 실패하였습니다.';
        }
    }
}

module.exports = new Vote();