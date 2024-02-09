import path from "path";
import fs from "fs";

// ensure directory exists
const configDir = "/etc/rancher-desktop"
fs.mkdirSync(configDir, { recursive: true });

// populate locked.json
const lockedPath = path.join(configDir, "locked.json");
const lockedObj = {
  version: 6,
  containerEngine: {
    allowedImages: {
      enabled: true,
      patterns: [
        "rancher/mirrored-[^/]+"
      ]
    }
  }
};
const lockedContents = JSON.stringify(lockedObj);
fs.writeFileSync(lockedPath, lockedContents);

// populate defaults.json
const defaultsPath = path.join(configDir, "defaults.json");
const defaultsObj = {
  version: 6,
  application: {
    autoStart: true
  },
  containerEngine: {
    name: "moby"
  },
  kubernetes: {
    version: "1.24.10"
  }
} 
const defaultsContents = JSON.stringify(defaultsObj);
fs.writeFileSync(defaultsPath, defaultsContents);
