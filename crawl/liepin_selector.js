const liepinSelector = {
    passwordLogin:
        "#user-reglpt > div.banner > div.wrap > div > div.login-box > ul > li:nth-child(2)",
    usernameInput: ".user-name",
    passwordInput:
        "#user-reglpt > div.banner > div.wrap > div > div.login-box > section.pc-username > div > form > div.verify-wrap > input",
    loginButton:
        "#user-reglpt > div.banner > div.wrap > div > div.login-box > section.pc-username > div > form > div.form-actions > button",
    logMsg: ".login-msg",
    findPeople: "#root > div.lpt-header > nav > div > ul > li:nth-child(3) > a",
    searchInput:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.search-bar > div.search-input-box > div > div > input",
    searchButton:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.search-bar > button.ant-btn.search-btn.ant-btn-primary",
    nextPage:
        "#root > div.wrap > div.board.search-resume-container.submited > div.resume-list-box > div.resume-list-pagebar > ul > li.ant-pagination-next > a",
    moreConditionSwitch:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.relative > a > span",
    randomClick:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.search-bar > div.search-input-box > div > div > input",
    inputGroup:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.search-bar > div.search-input-box > div > div",
    ageFrom:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.other-item.more > ul > li:nth-child(3) > div > dl.clearfix.age-container > dd > div > div > label:nth-child(1) > input",
    ageTo:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.other-item.more > ul > li:nth-child(3) > div > dl.clearfix.age-container > dd > div > div > label:nth-child(3) > input",
    本科:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.relative > ul > li:nth-child(4) > dl > dd > label:nth-child(2) > span",
    硕士:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.relative > ul > li:nth-child(4) > dl > dd > label:nth-child(3) > span",
    博士:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.relative > ul > li:nth-child(4) > dl > dd > label:nth-child(4) > span",
    大专:
        "#root > div.wrap > div.board.search-resume-container > div.filter-box > div.box-wrapper > div.conditons.relative > ul > li:nth-child(4) > dl > dd > label:nth-child(5) > span",
    antModalConfirm:
        "body > div:nth-child(19) > div > div.ant-modal-wrap > div > div.ant-modal-content > div > div > div.ant-modal-confirm-body"
};

module.exports = liepinSelector;
