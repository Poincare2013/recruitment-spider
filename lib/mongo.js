const MongoClient = require('mongodb').MongoClient;
const config = require('./config.js')
const url = config.mongoUrl;
const dbName = config.mongoDbName;

var ObjectId = require('mongodb').ObjectId;

class MyMgo {
    constructor() {
        this.client = new MongoClient(url, { useUnifiedTopology: true });
        this.db = null;
    }

    connect() {
        return new Promise((resolve, reject) => {
            this.client.connect((err, client) => {
                if (err) {
                    reject();
                } else {
                    console.log("connect success to mongo server");
                    this.db = client.db(dbName);
                    resolve();
                }
            })
        })
    }

    close() {
        this.client.close();
    }

    upsert(collectionName, data) {
        return new Promise((resolve, reject) => {
            let collection = this.db.collection(collectionName);
            let id = null;
            if (data._id == null) {
                id = new ObjectId();
            } else {
                id = data._id;
            }
            collection.updateOne({ '_id': id }, { $set: data }, { upsert: true }).then(result => {
                if (result.result.n === 1) {
                    // console.log("插入id: " + id + " 成功");
                    resolve(id);
                } else {
                    // console.log("插入id: " + id + " 失败");
                    reject(id);
                }
            }).catch(e => console.log(e));
        })
    }

}

module.exports = {
    MyMgo: MyMgo,
    dbName: dbName
}