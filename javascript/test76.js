import path from "path";
import fs from "fs";
import os from "os";

// ensure directory exists
const configDir = path.join(os.homedir(), ".config");
fs.mkdirSync(configDir, { recursive: true });

// populate locked.json
const lockedPath = path.join(configDir, "rancher-desktop.locked.json");
const lockedObj = {
  containerEngine: {
    allowedImages: {
      enabled: false
    }
  }
}
const lockedContents = JSON.stringify(lockedObj);
fs.writeFileSync(lockedPath, lockedContents);

// populate defaults.json
const defaultsPath = path.join(configDir, "rancher-desktop.defaults.json");
const defaultsObj = {
  version: 6,
  application: {
    autoStart: false
  },
  containerEngine: {
    name: "containerd"
  },
  kubernetes: {
    version: "1.23.16"
  }
} 
const defaultsContents = JSON.stringify(defaultsObj);
fs.writeFileSync(defaultsPath, defaultsContents);
