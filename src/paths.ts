import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";

export const ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "..");
export const CONFIG_DIR = path.join(os.homedir(), ".investiq");
export const CONFIG_PATH = path.join(CONFIG_DIR, "config.json");
