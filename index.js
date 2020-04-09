const puppeteer = require("puppeteer");
const rabbit = require("./lib/rabbit.js");
const mongo = require("./lib/mongo.js");
const file = require("./lib/file");
const utils = require("./lib/utils.js");
const liepinSelector = require("./crawl/liepin_selector.js");
const fs = require("fs");

const init_page = async function (page) {
    await page.setViewport({ width: 1280, height: 700 });
    await page.evaluateOnNewDocument(function () {
        Object.defineProperty(navigator, "webdriver", {
            get: function () {
                return false;
            },
        });
    });
};

const crawl = async function (browser, target, keyWord, education, ages) {
    let page = await browser.newPage();
    await init_page(page);

    let userAgent =
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36";
    await page.setUserAgent(userAgent);

    var data = fs
        .readFileSync("./go-recruitment-spider/config/liepin.json")
        .toString();

    var cookies = JSON.parse(data.replace(/(^\s*)|(\s*$)/g, ""));

    await page.setCookie(...cookies);

    var mgo = new mongo.MyMgo();
    const liepinCollection = "liepinCrawl";
    const liepinQueue = "liepin_list_page";

    await mgo.connect();
    const db = mgo.client.db(mongo.dbName);
    const collection = db.collection(liepinCollection);

    page.on("response", (response) => {
        //列表页
        if (
            response
                .url()
                .startsWith("https://lpt.liepin.com/cvsearch/search.json?")
        ) {
            response.text().then((response) => {
                console.log("已获取一页数据.....");
                resp = JSON.parse(response.replace(/(^\s*)|(\s*$)/g, ""));
                if (
                    resp &&
                    resp.data &&
                    resp.data.cvSearchResultForm &&
                    resp.data.cvSearchResultForm.cvSearchListFormList
                ) {
                    cvList = resp.data.cvSearchResultForm.cvSearchListFormList;
                    var resIdEncodes = [];
                    for (var item of cvList) {
                        //存mongo
                        collection
                            .updateOne(
                                { _id: item.resIdEncode },
                                { $addToSet: item },
                                { upsert: true }
                            )
                            .then(function () {
                                console.log("插入" + item.resIdEncode);
                            })
                            .catch(function (error) {
                                console.log(error);
                            });
                        //推mq
                        resIdEncodes.push(item.resIdEncode);
                    }
                    rabbit.send(liepinQueue, resIdEncodes);
                }
            });
        }
    });

    await page.goto(target);

    //登录部分
    // await page.waitFor(liepinSelector.passwordLogin);
    // await page.click(liepinSelector.passwordLogin);
    // await page.type(liepinSelector.usernameInput, "pg2020");
    // await page.type(liepinSelector.passwordInput, "pinguanzp");
    // await page.click(liepinSelector.loginButton);
    // await page.waitFor(liepinSelector.logMsg);

    await page.waitFor(liepinSelector.findPeople);

    await page.waitFor(1000);

    await page.click(liepinSelector.findPeople);

    await page.waitFor(liepinSelector.searchInput);

    //选择学历
    if (liepinSelector[education] == undefined) {
        console.log("未找到相应学历");
        browser.close();
    }
    await page.click(liepinSelector[education]);
    //输入关键字
    await page.type(liepinSelector.searchInput, keyWord);
    //选择年龄
    // await page.hover(liepinSelector.moreConditionSwitch);
    if (ages != undefined && ages.length == 2) {
        await page.click(liepinSelector.moreConditionSwitch);
        await page.waitFor(liepinSelector.ageFrom);
        await page.type(liepinSelector.ageFrom, ages[0]);
        await page.type(liepinSelector.ageTo, ages[1]);
    }
    await page.click(liepinSelector.searchButton);
    let curPage = 1;
    console.log("开始爬第" + curPage + "页");
    await page.waitFor(liepinSelector.nextPage);
    await page.hover(liepinSelector.nextPage);

    let ariaDisabled = await page.evaluate(function () {
        let next = document.getElementsByClassName("ant-pagination-next")[0];
        if (next) return next.getAttribute("aria-disabled");
        return "true";
    });

    //为true,可以点击下一页
    while (ariaDisabled == "false") {
        curPage++;
        console.log("开始爬第" + curPage + "页");
        await page.click(liepinSelector.nextPage);
        await page.waitFor(liepinSelector.nextPage);
        await page.hover(liepinSelector.nextPage);
        ariaDisabled = await page.evaluate(function () {
            let next = document.getElementsByClassName("ant-pagination-next");
            if (next && next.length > 0)
                return next[0].getAttribute("aria-disabled");
            return "true";
        });
        //等5秒
        await page.waitFor(5000);
        if (ariaDisabled == "true") {
            break;
        }
        //看看有没有封禁
        antModalConfirm = await page.evaluate(function () {
            let antModalConfirm = document.getElementsByClassName(
                "ant-modal-confirm-body"
            );
            if (antModalConfirm && antModalConfirm.length > 0) return true;
            return false;
        });

        if (antModalConfirm) {
            console.log("已封禁,结束.......");
            break;
        }
    }

    // await file.readSyncByRl(
    //     "# waiting validate, press enter after valid success:) >"
    // );

    page.close();
    mgo.client.close();
};

//main
(async () => {
    const browser = await puppeteer.launch({
        args: ["--window-size=1280,700"],
        timeout: 15000,
        ignoreHTTPSErrors: false,
        devtools: false,
        headless: false,
        // slowMo: 2500,
    });

    const keyWord = "工程师";
    const education = "本科"; //本科,大专,硕士,博士
    var ages = ["18", "39"];

    // let target = "https://passport.liepin.com/account/v1/elogin";

    let target = "https://lpt.liepin.com/cvsearch/showcondition/";

    await crawl(browser, target, keyWord, education, ages);

    browser.close();
})();
