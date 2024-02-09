import semver from 'npm:semver';

const versions = [
  'v2.7.0-plugins.beta.1',
  'v2.7.0-beta.1',
]

console.log(versions.sort(semver.rcompare));