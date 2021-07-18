const replaceInFile = require('replace-in-file');
const packageJson = require('./package.json');
const packageVersion = packageJson.version;

replaceInFile.sync({
    files: './src/environments/environment.prod.ts',
    from: [/version: '.+'/, /buildTimestamp: .+/],
    to: [`version: '${packageVersion}'`, `buildTimestamp: ${Math.floor(Date.now() / 1000)}`]
});
