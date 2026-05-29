import fs from "node:fs";
import { z } from "zod";
import { CONFIG_DIR, CONFIG_PATH } from "./paths.js";
import type { InvestiqConfig, PlatformConfig } from "./types.js";
import { hasCommand } from "./utils/shell.js";

export const platformSchema = z.object({
  label: z.string(),
  command: z.string(),
  args: z.array(z.string()),
  promptMode: z.enum(["argument", "print"]),
  subscription: z.string(),
  enabled: z.boolean()
});

export const configSchema = z.object({
  version: z.literal(1),
  defaultPlatform: z.string(),
  platformPriority: z.array(z.string()),
  requireFreshData: z.boolean(),
  allowTradeExecution: z.boolean(),
  platforms: z.record(z.string(), platformSchema)
});

export const PLATFORM_DEFAULTS: Record<string, PlatformConfig> = {
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

export function defaultConfig(): InvestiqConfig {
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

export function loadConfig(): InvestiqConfig {
  const fallback = defaultConfig();
  if (!fs.existsSync(CONFIG_PATH)) return fallback;

  const parsed = configSchema.partial().parse(JSON.parse(fs.readFileSync(CONFIG_PATH, "utf8")));
  return configSchema.parse({
    ...fallback,
    ...parsed,
    platforms: {
      ...fallback.platforms,
      ...parsed.platforms
    }
  });
}

export function saveConfig(config: InvestiqConfig): void {
  fs.mkdirSync(CONFIG_DIR, { recursive: true });
  fs.writeFileSync(CONFIG_PATH, `${JSON.stringify(config, null, 2)}\n`);
}
