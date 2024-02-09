import { spawn } from 'node:child_process';

interface SpawnError extends Error {
  command?: string[];
  stdout?: string;
  stderr?: string;
  code?: number;
  signal?: NodeJS.Signals;
}

interface SpawnOptionsEncoding {
  /**
   * The expected encoding of the executable.  If set, we will attempt to
   * convert the output to strings.
   */
  encoding?: { stdout?: BufferEncoding, stderr?: BufferEncoding } | BufferEncoding
}

/**
 * ErrorCommand is a symbol we attach to any exceptions thrown to describe the
 * command that failed.
 */
export const ErrorCommand = Symbol('child-process.command');

/**
 * Wrapper around child_process.spawn() to promisify it.  On Windows, we never
 * spawn a new command prompt window.
 * @param command The executable to spawn.
 * @param args Any arguments to the executable.
 * @param options Options to child_process.spawn();
 * @throws {SpawnError} When the command returns a failure
 */
export async function spawnFile(
  command: string,
): Promise<Record<string, never>>;
export async function spawnFile(
  command: string,
  options: SpawnOptionsWithStdioTuple<StdioNull | StdioPipe, StdioPipe, StdioPipe> & SpawnOptionsEncoding,
): Promise<{ stdout: string, stderr: string }>;
export async function spawnFile(
  command: string,
  args: string[],
): Promise<Record<string, never>>;
export async function spawnFile(
  command: string,
  args: string[],
  options: SpawnOptionsWithStdioTuple<StdioNull | StdioPipe, StdioPipe, StdioPipe> & SpawnOptionsEncoding,
): Promise<{ stdout: string, stderr: string }>;

/* eslint-enable no-redeclare */

export async function spawnFile(
  command: string,
  args?: string[],
  options: SpawnOptionsEncoding = {},
): Promise<{ stdout?: string, stderr?: string }> {
  let finalArgs: string[] = [];

  if (args && !Array.isArray(args)) {
    options = args;
    finalArgs = [];
  } else {
    finalArgs = args ?? [];
  }

  const stdio = options.stdio;
  const encodings = [
    undefined, // stdin
    (typeof options.encoding === 'string') ? options.encoding : options.encoding?.stdout,
    (typeof options.encoding === 'string') ? options.encoding : options.encoding?.stderr,
  ];
  const stdStreams: [stream.Readable | undefined, stream.Writable | undefined, stream.Writable | undefined] = [undefined, undefined, undefined];
  let mungedStdio: StdioOptions = 'pipe';

  // If we're piping to a stream, and we need to override the encoding, then
  // we need to do setup here.
  if (Array.isArray(stdio)) {
    mungedStdio = ['ignore', 'ignore', 'ignore'];
    for (let i = 0; i < 3; i++) {
      const original = stdio[i];
      let munged: StdioNull | StdioPipe | number;

      if (i === 0 && original instanceof stream.Readable) {
        munged = 'pipe';
        stdStreams[i] = original;
      } else {
        munged = original;
      }
      if (munged instanceof stream.Writable && encodings[i]) {
        stdStreams[i] = munged;
        munged = 'pipe';
      }
      mungedStdio[i] = munged;
    }
  } else if (typeof stdio === 'string') {
    mungedStdio = [stdio, stdio, stdio];
  }

  // Spawn the child, overriding options.stdio.  This is necessary to support
  // transcoding the output.
  const child = spawn(command, finalArgs, {
    windowsHide: true,
    ...options,
    stdio:       mungedStdio,
  });
  const resultMap: Record<number, 'stdout' | 'stderr'> = { 1: 'stdout', 2: 'stderr' };
  const result: { stdout?: string, stderr?: string } = {};

  if (Array.isArray(mungedStdio)) {
    if (stdStreams[0] instanceof stream.Readable && child.stdin) {
      stdStreams[0].pipe(child.stdin);
    }
    for (const i of [1, 2] as const) {
      if (mungedStdio[i] === 'pipe') {
        const encoding = encodings[i];
        const childStream = child[resultMap[i]];

        if (!stdStreams[i]) {
          result[resultMap[i]] = '';
        }
        if (childStream) {
          if (encoding) {
            childStream.setEncoding(encoding);
          }
          childStream.on('data', (chunk) => {
            if (stdStreams[i]) {
              (stdStreams[i] as stream.Writable).write(chunk);
            } else {
              result[resultMap[i]] += chunk;
            }
          });
        }
      }
    }
  }

  await new Promise<void>((resolve, reject) => {
    child.on('exit', (code, signal) => {
      if ((code === 0 && signal === null) || (code === null && signal === 'SIGTERM')) {
        return resolve();
      }
      let message = `${ command } exited with code ${ code }`;

      if (code === null) {
        message = `${ command } exited with signal ${ signal }`;
      }
      const error: SpawnError = new Error(message);

      Object.defineProperties(error, {
        [ErrorCommand]: {
          enumerable: false,
          value:      `${ command } ${ finalArgs.join(' ') }`,
        },
      });
      if (typeof result.stdout !== 'undefined') {
        error.stdout = result.stdout;
      }
      if (typeof result.stderr !== 'undefined') {
        error.stderr = result.stderr;
      }
      if (code !== null) {
        error.code = code;
      } else if (signal !== null) {
        error.signal = signal;
      }
      error.command = [command].concat(finalArgs);
      reject(error);
    });
    child.on('error', reject);
  });

  return result;
}

const result = await spawnFile('false');
