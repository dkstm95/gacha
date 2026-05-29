import { spawnSync } from "node:child_process";

export function shellQuote(value: string): string {
  return `'${String(value).replaceAll("'", "'\\''")}'`;
}

export function hasCommand(command: string): boolean {
  if (!command) return false;
  const result = spawnSync("sh", ["-lc", `command -v ${shellQuote(command)}`], {
    encoding: "utf8"
  });
  return result.status === 0;
}
