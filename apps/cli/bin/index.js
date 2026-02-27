#!/usr/bin/env node

import { spawn } from 'node:child_process';
import fs from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Mapping OS and architecture to the folders created by GoReleaser
const platforms = {
  darwin: {
    x64: 'darwin_amd64_v1/dropdx',
    arm64: 'darwin_arm64/dropdx',
  },
  linux: {
    x64: 'linux_amd64_v1/dropdx',
    arm64: 'linux_arm64/dropdx',
  },
  win32: {
    x64: 'windows_amd64_v1/dropdx.exe',
    arm64: 'windows_arm64/dropdx.exe',
  },
};

const platform = os.platform();
const arch = os.arch();

let binaryPath = '';

// Check if we are in development mode (root binary)
const rootBinary = path.join(__dirname, '../../../dropdx');
if (fs.existsSync(rootBinary)) {
  binaryPath = rootBinary;
} else if (platforms[platform] && platforms[platform][arch]) {
  binaryPath = path.join(__dirname, platforms[platform][arch]);
}

if (!binaryPath || !fs.existsSync(binaryPath)) {
  console.error(
    `Unsupported platform/architecture or binary not found: ${platform}/${arch}`
  );
  process.exit(1);
}

// Ensure the binary is executable
try {
  fs.chmodSync(binaryPath, 0o755);
} catch {
  // Ignore error if we don't have permission to chmod (it might be read-only filesystem)
  // but if it's already executable it will work anyway.
}

const args = process.argv.slice(2);

const child = spawn(binaryPath, args, { stdio: 'inherit' });

child.on('close', (code) => {
  process.exit(code);
});
