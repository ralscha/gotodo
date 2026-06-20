const { replaceInFileSync } = require('replace-in-file');
const packageJson = require('./package.json');
const packageVersion = packageJson.version;

replaceInFileSync({
  files: './src/environments/environment.prod.ts',
  from: [/version: '.+'/, /buildTimestamp: .+/],
  to: [`version: '${packageVersion}'`, `buildTimestamp: ${Math.floor(Date.now() / 1000)}`],
});
