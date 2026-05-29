export const MODES = ["discover", "select", "entry", "exit", "portfolio", "journal"] as const;

export type Mode = (typeof MODES)[number];
export type PlatformName = "claude" | "codex" | "opencode" | "gemini" | "manual" | string;
export type PromptMode = "argument" | "print";

export interface PlatformConfig {
  label: string;
  command: string;
  args: string[];
  promptMode: PromptMode;
  subscription: string;
  enabled: boolean;
}

export interface InvestiqConfig {
  version: 1;
  defaultPlatform: "auto" | PlatformName;
  platformPriority: string[];
  requireFreshData: boolean;
  allowTradeExecution: boolean;
  platforms: Record<string, PlatformConfig>;
}

export interface RunOptions {
  platform: string;
  dryRun: boolean;
  query: string[];
}
