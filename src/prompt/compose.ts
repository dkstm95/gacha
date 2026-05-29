import fs from "node:fs";
import path from "node:path";
import { ROOT } from "../paths.js";
import { MODES, type Mode } from "../types.js";

function readText(relativePath: string): string {
  return fs.readFileSync(path.join(ROOT, relativePath), "utf8");
}

export function isMode(value: string): value is Mode {
  return (MODES as readonly string[]).includes(value);
}

export function buildPrompt(mode: Mode, queryParts: string[]): string {
  const query = queryParts.join(" ").trim();
  const system = readText("plugins/investiq/platforms/generic/system-prompt.md");
  const template = readText("plugins/investiq/templates/investment-report.md");
  const workflowPath = `plugins/investiq/workflows/${mode}.md`;
  const absoluteWorkflowPath = path.join(ROOT, workflowPath);
  const workflow = fs.existsSync(absoluteWorkflowPath)
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
