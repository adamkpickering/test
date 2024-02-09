import child_process from 'child_process';

try {
  const result = child_process.spawnSync('aewrsdf', ['doesnotexist.txt'], {encoding: 'utf8'});
  console.log(JSON.stringify(result)); 
} catch (error) {
  console.log(`Got error ${error}`);
}
