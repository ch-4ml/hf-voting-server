'use strict';

const express = require('express');
const router = express.Router();
const candidateModel = require('../model/candidate');

router.post('/new', async(req, res) => {
    const name = req.body.name;
    const voteID = req.body.voteID;
    let data;
    try {
        const result = await candidateModel.setCandidate(name, voteID);
        data = { msg: result, result: true }
        res.status(200).json({ data: data })
    } catch(err) {
        data = { msg: err, result: false }
        res.status(500).json({ data: data });
    }
});

module.exports = router;