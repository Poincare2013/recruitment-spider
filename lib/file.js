const readline = require("readline");
const fs = require("fs");

function readSyncByRl(tips) {
    tips = tips || "> ";
    return new Promise((resolve) => {
        const rl = readline.createInterface({
            input: process.stdin,
            //output: process.stdout
            output: process.stderr,
        });
        rl.question(tips, (answer) => {
            rl.close();
            resolve(answer.trim());
        });
    });
}

function readFileLines(filename) {
    return new Promise((resolve) => {
        let lines = [];
        const rl = readline.createInterface({
            input: fs.createReadStream(filename),
        });
        rl.on("line", (line) => {
            lines.push(line.split(" ", 1)[0]);
        });
        rl.on("close", () => {
            rl.close();
            resolve(lines);
        });
    });
}

module.exports = {
    readSyncByRl: readSyncByRl,
    readFileLines: readFileLines,
};
