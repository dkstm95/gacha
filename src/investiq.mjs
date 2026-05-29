#!/usr/bin/env node
import fs from "node:fs";
import os from "node:os";
import path from "node:path";
import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";
import { spawnSync } from "node:child_process";

const ROOT = path.resolve(path.dirname(new URL(import.meta.url).pathname), "..");
const CONFIG_DIR = path.join(os.homedir(), ".investiq");
const CONFIG_PATH = path.join(CONFIG_DIR, "config.json");
const MODES = new Set(["discover", "select", "entry", "exit", "portfolio", "journal"]);

const PLATFORM_DEFAULTS = {
  claude: {
    label: "Claude Code",
    command: "claude",
    args: ["-p", "{{prompt}}"],
    promptMode: "argument",
    subscription: "",
    enabled: false
  },
  codex: {
    label: "Codex",
    command: "codex",
    args: ["{{prompt}}"],
    promptMode: "argument",
    subscription: "",
    enabled: false
  },
  opencode: {
    label: "OpenCode / Oh My OpenAgent",
    command: "opencode",
    args: ["run", "{{prompt}}"],
    promptMode: "argument",
    subscription: "",
    enabled: false
  },
  gemini: {
    label: "Gemini CLI",
    command: "gemini",
    args: ["{{prompt}}"],
    promptMode: "argument",
    subscription: "",
    enabled: false
  },
  manual: {
    label: "Manual copy/paste",
    command: "",
    args: [],
    promptMode: "print",
    subscription: "manual",
    enabled: true
  }
};

function readText(relativePath) {
  return fs.readFileSync(path.join(ROOT, relativePath), "utf8");
}

function hasCommand(command) {
  if (!command) return false;
  const result = spawnSync("sh", ["-lc", `command -v ${shellQuote(command)}`], {
    encoding: "utf8"
  });
  return result.status === 0;
}

function shellQuote(value) {
  return `'${String(value).replaceAll("'", "'\\''")}'`;
}

function defaultConfig() {
  const platforms = structuredClone(PLATFORM_DEFAULTS);
  for (const platform of Object.values(platforms)) {
    if (platform.command && hasCommand(platform.command)) {
      platform.enabled = true;
    }
  }
  return {
    version: 1,
    defaultPlatform: "auto",
    platformPriority: ["claude", "codex", "opencode", "gemini", "manual"],
    requireFreshData: true,
    allowTradeExecution: false,
    platforms
  };
}

function loadConfig() {
  if (!fs.existsSync(CONFIG_PATH)) return defaultConfig();
  return {
    ...defaultConfig(),
    ...JSON.parse(fs.readFileSync(CONFIG_PATH, "utf8"))
  };
}

function saveConfig(config) {
  fs.mkdirSync(CONFIG_DIR, { recursive: true });
  fs.writeFileSync(CONFIG_PATH, `${JSON.stringify(config, null, 2)}\n`);
}

function buildPrompt(mode, queryParts) {
  if (!MODES.has(mode)) {
    throw new Error(`Unknown mode: ${mode}`);
  }

  const query = queryParts.join(" ").trim();
  const system = readText("plugins/investiq/platforms/generic/system-prompt.md");
  const template = readText("plugins/investiq/templates/investment-report.md");
  const workflowPath = `plugins/investiq/workflows/${mode}.md`;
  const workflow = fs.existsSync(path.join(ROOT, workflowPath))
    ? readText(workflowPath)
    : `# investiq ${mode}\n\nFollow the investiq generic workflow for ${mode}.`;

  return [
    system.trim(),
    "",
    workflow.trim(),
    "",
    "User request:",
    query || "(No additional user request supplied.)",
    "",
    "Report template:",
    template.trim(),
    "",
    "Hard requirements:",
    "- Use current web search or current market-data tools before analysis.",
    "- If fresh data cannot be verified, do not make a recommendation.",
    "- Include data freshness, source links, risks, Devil's Advocate, action conditions, monitoring plan, and provenance.",
    "- Do not execute trades. The final decision remains with the user."
  ].join("\n");
}

function selectPlatform(config, requested) {
  if (requested && requested !== "auto") {
    return requested;
  }
  for (const name of config.platformPriority) {
    const platform = config.platforms[name];
    if (!platform?.enabled) continue;
    if (name === "manual" || hasCommand(platform.command)) return name;
  }
  return "manual";
}

function renderArgs(args, prompt) {
  return args.map((arg) => arg.replaceAll("{{prompt}}", prompt));
}

function runPlatform(name, platform, prompt, dryRun) {
  if (name === "manual" || platform.promptMode === "print" || !platform.command) {
    console.log(prompt);
    return 0;
  }

  if (!hasCommand(platform.command)) {
    console.error(`Platform command not found: ${platform.command}`);
    console.error("Run `investiq doctor` or use `--platform manual`.");
    return 1;
  }

  const args = renderArgs(platform.args || [], prompt);
  if (dryRun) {
    console.log([platform.command, ...args.map(shellQuote)].join(" "));
    return 0;
  }

  const result = spawnSync(platform.command, args, {
    stdio: "inherit",
    env: process.env
  });
  return result.status ?? 1;
}

async function initConfig(args) {
  const yes = args.includes("--yes") || args.includes("-y");
  const config = defaultConfig();

  if (!yes) {
    const rl = readline.createInterface({ input, output });
    console.log("Configure the AI platforms you subscribe to or can run locally.");
    for (const [name, platform] of Object.entries(config.platforms)) {
      if (name === "manual") continue;
      const installed = hasCommand(platform.command);
      const answer = await rl.question(
        `Enable ${platform.label} (${platform.command})? ${installed ? "[Y/n]" : "[y/N]"} `
      );
      const normalized = answer.trim().toLowerCase();
      platform.enabled = installed
        ? normalized !== "n" && normalized !== "no"
        : normalized === "y" || normalized === "yes";
      if (platform.enabled) {
        const sub = await rl.question(`Subscription/account label for ${platform.label} (optional): `);
        platform.subscription = sub.trim();
      }
    }
    await rl.close();
  }

  saveConfig(config);
  console.log(`Wrote ${CONFIG_PATH}`);
}

function doctor() {
  const config = loadConfig();
  console.log(`Config: ${fs.existsSync(CONFIG_PATH) ? CONFIG_PATH : "(not created; using defaults)"}`);
  console.log("");
  for (const [name, platform] of Object.entries(config.platforms)) {
    const installed = name === "manual" ? true : hasCommand(platform.command);
    const status = platform.enabled && installed ? "ready" : platform.enabled ? "missing" : "disabled";
    console.log(`${name.padEnd(9)} ${status.padEnd(8)} ${platform.label}`);
    if (platform.command) console.log(`           command: ${platform.command}`);
    if (platform.subscription) console.log(`           subscription: ${platform.subscription}`);
  }
}

function listPlatforms() {
  const config = loadConfig();
  console.log(JSON.stringify(config.platforms, null, 2));
}

function usage() {
  console.log(`investiq

Usage:
  investiq init [--yes]                       Create ~/.investiq/config.json
  investiq doctor                             Check configured AI platforms
  investiq platforms                          Print platform config
  investiq prompt <mode> [request]            Print the composed agent prompt
  investiq run <mode> [request] [options]     Run through a configured platform
  investiq <mode> [request]                   Alias for prompt <mode>

Modes:
  discover, select, entry, exit, portfolio, journal

Options:
  --platform auto|claude|codex|opencode|gemini|manual
  --dry-run

Examples:
  investiq init
  investiq doctor
  investiq run entry "NVDA current entry zone" --platform auto
  investiq run discover "latest opportunities for a 12 month horizon" --platform manual
`);
}

async function main(argv) {
  const [command, ...rest] = argv;

  if (!command || command === "help" || command === "-h" || command === "--help") {
    usage();
    return 0;
  }

  if (command === "init") {
    await initConfig(rest);
    return 0;
  }

  if (command === "doctor") {
    doctor();
    return 0;
  }

  if (command === "platforms") {
    listPlatforms();
    return 0;
  }

  if (command === "prompt" || MODES.has(command)) {
    const mode = command === "prompt" ? rest.shift() : command;
    if (!mode || !MODES.has(mode)) throw new Error("Missing or invalid mode.");
    console.log(buildPrompt(mode, rest));
    return 0;
  }

  if (command === "run") {
    const mode = rest.shift();
    if (!mode || !MODES.has(mode)) throw new Error("Missing or invalid mode.");
    const dryRun = rest.includes("--dry-run");
    const platformIndex = rest.indexOf("--platform");
    const requested = platformIndex >= 0 ? rest[platformIndex + 1] : "auto";
    const query = rest.filter((arg, index) => {
      if (arg === "--dry-run") return false;
      if (arg === "--platform") return false;
      if (platformIndex >= 0 && index === platformIndex + 1) return false;
      return true;
    });
    const config = loadConfig();
    const platformName = selectPlatform(config, requested);
    const platform = config.platforms[platformName];
    if (!platform) throw new Error(`Unknown platform: ${platformName}`);
    return runPlatform(platformName, platform, buildPrompt(mode, query), dryRun);
  }

  throw new Error(`Unknown command: ${command}`);
}

main(process.argv.slice(2))
  .then((code) => process.exit(code))
  .catch((error) => {
    console.error(error.message);
    process.exit(1);
  });
