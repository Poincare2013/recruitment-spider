const config = {
    rabbitUrl: 'amqp://admin:admin@localhost:5672',
    mongoUrl: 'mongodb://crawl:crawl@localhost:27017?authSource=crawl',
    mongoDbName: 'crawl',
    liepinListPageQueue: 'liepin_list_page',
}

module.exports = config;