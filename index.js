const express = require('express');
const compression = require('compression');
const line = require("@line/bot-sdk");
const dotenv = require("dotenv");
const helmet = require("helmet");

dotenv.config();

const config = {
	channelAccessToken: process.env.LINE_CHANNEL_ACCESS_TOKEN,
	channelSecret: process.env.LINE_CHANNEL_SECRET,
};
const client = new line.Client(config);
const port = process.env.PORT || 3000;

const app = express();
// https://c2f4d665da7c.ngrok.io/

app.use(helmet());
app.use(compression())

const handleLineEvent = (req, res) => {
    const event = req.body.events[0]
	console.log(event)
    let reply = {
        type: "text",
        text: event.message.text,
    };
	if (event.type !== "message") {
        res.send('OK but not ok');
	}

	if (event.source.type !== "user") {
        res.send('OK but not ok');
	}
    client.replyMessage(event.replyToken, reply)
    res.send('OK');
};

app.get('/', (req, res) => {
  res.send('Arjuna!')
})

app.post("/webhook", line.middleware(config), handleLineEvent);

app.listen(port, () => {
  console.log(`listening at port ${port}`)
})
