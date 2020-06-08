/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');

const ccpPath = path.resolve(__dirname, '..', 'first-network', 'connection-org1.json');

module.exports = async (id) => {
    try {
        console.log(ccpPath);
        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(process.cwd(), 'wallet');
        console.log(`walletPath: ${walletPath}`);
        const wallet = new FileSystemWallet(walletPath);

        // Check to see if we've already enrolled the user.
        const userExists = await wallet.exists(id);
        if (!userExists) {
            console.log(`An identity for the user ${id} does not exist in the wallet`);
            console.log(`Run the registerUser.js application before retrying`);
            return;
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        console.log(`gateway: ${gateway.toString()}`);
        const r = await gateway.connect(ccpPath, { wallet, identity: id, discovery: { enabled: true, asLocalhost: true } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('hf-voting');

        // Evaluate the specified transaction.
        const result = await contract.evaluateTransaction('getAllVotes');

        return result.toString();

    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        // process.exit(1);
    }
}