import semver from 'npm:semver';

const isoVersionRegex = /^[0-9]+\.[0-9]+\.[0-9]+\.rd[0-9]+$/;

function rcompareVersions(isoVersion1: string, isoVersion2: string): -1 | 0 | 1 {
  const isoVersion1Parts = isoVersion1.split('.');
  const isoVersion2Parts = isoVersion2.split('.');

  if (!isoVersionRegex.test(isoVersion1) || !isoVersionRegex.test(isoVersion2)) {
    throw new Error(`One or both of ${isoVersion1} and ${isoVersion2} are not in expected format ${isoVersionRegex}`);
  }

  const normalVersion1 = isoVersion1Parts.slice(0, 3).join('.');
  const normalVersion2 = isoVersion2Parts.slice(0, 3).join('.');
  const prereleaseVersion1 = isoVersion1Parts[3].replace('rd', '');
  const prereleaseVersion2 = isoVersion2Parts[3].replace('rd', '');
  const semverVersion1 = `${normalVersion1}-${prereleaseVersion1}`;
  const semverVersion2 = `${normalVersion2}-${prereleaseVersion2}`;

  return semver.rcompare(semverVersion1, semverVersion2);
}

const versions = [
  '0.1.0.rd6',
  '0.1.0.rd5',
  '0.1.0.rd10',
  '0.1.0.rd2',
  '0.2.0.rd10',
  '0.2.0.rd11',
  '0.3.1.rd1',
]

console.log(versions.sort(rcompareVersions));
