import semver from 'npm:semver';

const inputs = [
  'v3.1',
  '3.1',
  '  3.4',
  ' v3.4.1',
  '2.3.1  ',
  '2.3  ',
  'desktop-v2.7.0.beta.1',
  'desktop-v2.7.0-beta.1',
  '2.7.0-beta.1',
  'v2.7.0-beta.1',
]

for (const input of inputs) {
  console.log(`input: "${input}"\toutput: ${semver.coerce(input)}\t input valid: ${semver.valid(input)}`);
}