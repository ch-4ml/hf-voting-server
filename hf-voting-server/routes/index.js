'use strict';

const express = require('express');
const router = express.Router();
const crypto = require('crypto');
const fs = require('fs');
const nodeZip = require('node-zip');
const unzip = require('unzip');
const path = require('path');
const register = require('../registerUser');
const query = require('../query');
const multer = require('multer');
const upload = multer.diskStorage({
    destination: function(req, file, cb) {
        cb(null, 'wallet/')
    },
    filename: function(req, file, cb) {
        cb(null, file.originalname)
    }
})
const wallet = multer({storage: upload});
const walletPath = path.join(__dirname, '/../wallet');

router.get('/', (req, res) => {
    res.render('index');
});

router.get('/register', async (req, res) => {
    try {
        const id = await crypto.randomBytes(32).toString('hex')
        await register(id);
        
        // Make wallet file
        const files = fs.readdirSync(`${walletPath}/${id}`, { withFileTypes: true });

        // Generate zip object
        const zip = new nodeZip();

        for(let i = 0; i < files.length; i++) {
            zip.file(`${files[i]}`, fs.readFileSync(`${walletPath}/${id}/${files[i]}`));
        }

        const zipData = zip.generate({ base64: false, compression: 'DEFLATE' });
        fs.writeFileSync(`${walletPath}/${id}.zip`, zipData, 'binary');

        res.download(`${walletPath}/${id}.zip`, `${id}.zip`);

        // Remove wallet file
        for(let i = 0; i < files.length; i++) {
            fs.unlinkSync(`${walletPath}/${id}/${files[i]}`);
        }
        fs.rmdirSync(`${walletPath}/${id}`, { recursive: true });

    } catch(err) {
        console.log(err);
    }
});

router.post('/query', wallet.single('file'), async (req, res) => {
    // const { file: { path }} = req;
    try {
        const filename = req.file.filename;
        const id = filename.split('.')[0];
        fs.createReadStream(`${req.file.path}`).pipe(unzip.Extract({ path: `${walletPath}/${id}` }));
        const result = await query(id);
        console.log(result);
        res.json(JSON.parse(result));
    } catch (err) {
        console.log(err);
        res.json('Error occured!!');
    }
});

module.exports = router;
