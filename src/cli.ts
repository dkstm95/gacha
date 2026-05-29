#!/usr/bin/env node
import fs from "node:fs";
import readline from "node:readline/promises";
import { stdin as input, stdout as output } from "node:process";
import { Command } from "commander";
import { CONFIG_PATH } from "./paths.js";
import { defaultConfig, loadConfig, saveConfig } from "./config.js";
import { buildPrompt, isMode } from "./prompt/compose.js";
import { runPlatform, selectPlatform } from "./platforms/runner.js";
import { MODES, type Mode } from "./types.js";
import { hasCommand } from "./utils/shell.js";

export async function initConfig(yes: boolean): Promise<void> {
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

export function doctor(): void {
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

function parseQueryArgs(values: string[]): string[] {
  return values.filter((value) => value.length > 0);
}

export async function runCli(argv: string[]): Promise<number> {
  const program = new Command();

  program
    .name("investiq")
    .description("Fresh-data investment research agent harness")
    .version("0.1.0");

  program
    .command("init")
    .description("Create ~/.investiq/config.json")
    .option("-y, --yes", "accept detected platform defaults")
    .action(async (options: { yes?: boolean }) => {
      await initConfig(Boolean(options.yes));
    });

  program
    .command("doctor")
    .description("Check configured AI platforms")
    .action(() => doctor());

  program
    .command("platforms")
    .description("Print platform config")
    .action(() => {
      console.log(JSON.stringify(loadConfig().platforms, null, 2));
    });

  program
    .command("prompt")
    .description("Print the composed agent prompt")
    .argument("<mode>", `one of: ${MODES.join(", ")}`)
    .argument("[request...]", "investment request")
    .action((mode: string, request: string[] = []) => {
      if (!isMode(mode)) throw new Error(`Invalid mode: ${mode}`);
      console.log(buildPrompt(mode, parseQueryArgs(request)));
    });

  program
    .command("run")
    .description("Run through a configured platform")
    .argument("<mode>", `one of: ${MODES.join(", ")}`)
    .argument("[request...]", "investment request")
    .option("--platform <platform>", "auto|claude|codex|opencode|gemini|manual", "auto")
    .option("--dry-run", "print the platform command without executing it")
    .action((mode: string, request: string[] = [], options: { platform: string; dryRun?: boolean }) => {
      if (!isMode(mode)) throw new Error(`Invalid mode: ${mode}`);
      const config = loadConfig();
      const platformName = selectPlatform(config, options.platform);
      const platform = config.platforms[platformName];
      if (!platform) throw new Error(`Unknown platform: ${platformName}`);
      const prompt = buildPrompt(mode, parseQueryArgs(request));
      process.exitCode = runPlatform(platformName, platform, prompt, Boolean(options.dryRun));
    });

  for (const mode of MODES) {
    program
      .command(mode)
      .description(`Alias for: investiq prompt ${mode}`)
      .argument("[request...]", "investment request")
      .action((request: string[] = []) => {
        console.log(buildPrompt(mode as Mode, parseQueryArgs(request)));
      });
  }

  await program.parseAsync(argv, { from: "user" });
  return typeof process.exitCode === "number" ? process.exitCode : 0;
}
