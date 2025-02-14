import {replaceInFileSync} from 'replace-in-file'
import info from './package.json' with { type: "json" };
const packageVersion = info.version;
console.log(packageVersion);
replaceInFileSync({
    files: './src/environments/environment.prod.ts',
    from: [/version: '.+'/, /buildTimestamp: .+/],
    to: [`version: '${packageVersion}'`, `buildTimestamp: ${Math.floor(Date.now() / 1000)}`]
});
