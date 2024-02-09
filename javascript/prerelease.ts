import semver from 'npm:semver';

const versions = [
  '3.11.0-rc.1',
  '3.10.3',
  '3.10.2',
  '3.10.2-rc.2',
]

console.log('All versions:');
console.log(versions);
console.log('Non-prerelease versions:');
console.log(versions.filter(version => semver.prerelease(version) === null));
