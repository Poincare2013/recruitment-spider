const config = require('./config.js');
const url = config.rabbitUrl;
const amqp = require('amqplib');

function send(queue, msgs) {
	return new Promise((resolve, reject) => {
		// var amqp = require('amqplib');
		amqp.connect(url).then(function (conn) {
			return conn.createChannel().then(function (ch) {
				var ok = ch.assertQueue(queue, { durable: true });
				return ok.then(function () {
					for (var i = 0; i <= msgs.length - 1; i++) {
						ch.sendToQueue(queue, Buffer.from(msgs[i]), { deliveryMode: true });
						console.log(" [x] Sent '%s'", msgs[i]);
					}
					resolve()
					return ch.close();
				});
			}).finally(function () { conn.close() });
		}).catch(console.warn);
	})
}

//promise版本,一个一个接收,这里超级坑,只能get,还得是promise的
function receive(queue) {
	return new Promise((resolve, reject) => {
		amqp.connect(url).then(function (conn) {
			return conn.createChannel().then(function (ch) {
				var ok=ch.assertQueue(queue, { durable: true });
				return ok.then(function () {
					ch.get(queue, { noAck: true }).then(msg => {
						if (msg) {
							// console.log(" [x] Received '%s'", msg.content.toString());
							// ch.ack(msg)
							resolve(msg.content.toString())
						}else{
							resolve("over")
						}
						ch.close();
						return conn.close();
					})
				})
			});
		});
	})
}

(async function () {
	// await send(config.tmallRateQueue, ['gege', 'aaaa']);
	// let msg = await receive(config.tmallRateQueue);
	// console.log(msg);
	// let msg1 = await receive(config.tmallRateQueue);
	// console.log(msg1);
})()

module.exports.send = send
module.exports.receive = receive