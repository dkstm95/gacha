import { spawnSync } from "node:child_process";
import type { InvestiqConfig, PlatformConfig } from "../types.js";
import { hasCommand, shellQuote } from "../utils/shell.js";

export function selectPlatform(config: InvestiqConfig, requested: string): string {
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

export function renderArgs(args: string[], prompt: string): string[] {
  return args.map((arg) => arg.replaceAll("{{prompt}}", prompt));
}

export function runPlatform(name: string, platform: PlatformConfig, prompt: string, dryRun: boolean): number {
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
