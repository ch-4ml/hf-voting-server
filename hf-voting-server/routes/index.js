'use strict';

const express = require('express');
const router = express.Router();
const crypto = require('crypto');
const fs = require('fs');
const nodeZip = require('node-zip');
const register = require('../registerUser');

router.get('/', (req, res) => {
    res.render('index');
});

router.get('/register', async (req, res) => {
    try {
        const id = await crypto.randomBytes(32).toString('hex')
        await register(id);
        
        // Make wallet file
        const files = fs.readdirSync(`/tmp/wallet/${id}`, { withFileTypes: true });

        // Generate zip object
        const zip = new nodeZip();

        for(let i = 0; i < files.length; i++) {
            zip.file(`${files[i]}`, fs.readFileSync(`/tmp/wallet/${id}/${files[i]}`));
        }

        const zipData = zip.generate({ base64: false, compression: 'DEFLATE' });
        fs.writeFileSync(`/tmp/wallet/${id}.zip`, zipData, 'binary');

        res.download(`/tmp/wallet/${id}.zip`, `${id}.zip`);

        // Remove wallet file
        for(let i = 0; i < files.length; i++) {
            // fs.unlinkSync(`/tmp/wallet/${id}/${files[i]}`);
        }
        // fs.rmdirSync(`/tmp/wallet/${id}`, { recursive: true });

    } catch(err) {
        console.log(err);
        res.status(500).send(err);
    }
});

module.exports = router;
