#!/usr/bin/env node
import fs from "node:fs";
import path from "node:path";

const root = process.cwd();
const pluginRoot = path.join(root, "plugins", "investiq");
const manifestPath = path.join(pluginRoot, ".codex-plugin", "plugin.json");
const marketplacePath = path.join(root, ".agents", "plugins", "marketplace.json");

function fail(message) {
  console.error(message);
  process.exit(1);
}

function readJson(filePath) {
  try {
    return JSON.parse(fs.readFileSync(filePath, "utf8"));
  } catch (error) {
    fail(`Invalid JSON: ${filePath}\n${error instanceof Error ? error.message : String(error)}`);
  }
}

function requireFile(relativePath) {
  const absolutePath = path.join(root, relativePath);
  if (!fs.existsSync(absolutePath)) {
    fail(`Missing required file: ${relativePath}`);
  }
}

function requireString(object, field) {
  if (typeof object[field] !== "string" || object[field].trim() === "") {
    fail(`Missing required string field: ${field}`);
  }
}

requireFile("plugins/investiq/.codex-plugin/plugin.json");
requireFile(".agents/plugins/marketplace.json");
requireFile("plugins/investiq/skills/investiq/SKILL.md");
requireFile("plugins/investiq/templates/investment-report.md");

for (const mode of ["discover", "select", "entry", "exit", "portfolio", "journal"]) {
  requireFile(`plugins/investiq/workflows/${mode}.md`);
}

const manifest = readJson(manifestPath);
requireString(manifest, "name");
requireString(manifest, "version");
requireString(manifest, "description");

if (manifest.name !== "investiq") {
  fail(`Expected manifest name "investiq", got "${manifest.name}"`);
}

if (!/^\d+\.\d+\.\d+$/.test(manifest.version)) {
  fail(`Manifest version must be strict semver, got "${manifest.version}"`);
}

if (manifest.skills !== "./skills/") {
  fail('Manifest skills path must be "./skills/"');
}

for (const field of ["displayName", "shortDescription", "longDescription", "developerName", "category"]) {
  requireString(manifest.interface ?? {}, field);
}

const marketplace = readJson(marketplacePath);
const entry = marketplace.plugins?.find((plugin) => plugin.name === "investiq");
if (!entry) {
  fail("Marketplace does not include investiq plugin entry");
}

if (entry.source?.path !== "./plugins/investiq") {
  fail('Marketplace source.path must be "./plugins/investiq"');
}

if (!entry.policy?.installation || !entry.policy?.authentication || !entry.category) {
  fail("Marketplace entry must include policy.installation, policy.authentication, and category");
}

console.log(`Plugin validation passed: ${pluginRoot}`);
