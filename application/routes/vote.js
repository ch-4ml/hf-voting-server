'use strict';

const express = require('express');
const router = express.Router();
const voteModel = require('../model/vote');

router.get('/', async(req, res) => {
    const id = req.body.id;
    let data;
    try {
        const result = await voteModel.getVote(id);
        data = { msg: result, result: true }
        res.status(200).json({ data: data })
    } catch(err) {
        data = { msg: err, result: false }
        res.status(500).json({ data: data });
    }
});

router.get('/all', async(req, res) => {
    let data;
    try {
        const result = await voteModel.getAllVotes();
        data = { msg: result, result: true }
        res.status(200).json({ data: data })
    } catch(err) {
        data = { msg: err, result: false }
        res.status(500).json({ data: data });
    }
});

router.post('/new', async(req, res) => {
    const name = req.body.name;
    let data;
    try {
        const result = await voteModel.setVote(name);
        data = { msg: result, result: true }
        res.status(200).json({ data: data })
    } catch(err) {
        data = { msg: err, result: false }
        res.status(500).json({ data: data });
    }
});

module.exports = router;