import { describe, expect, it } from "vitest";
import { buildPrompt, isMode } from "../src/prompt/compose.js";

describe("prompt composition", () => {
  it("builds an entry prompt with hard requirements", () => {
    const prompt = buildPrompt("entry", ["NVDA"]);

    expect(prompt).toContain("# investiq entry");
    expect(prompt).toContain("User request:\nNVDA");
    expect(prompt).toContain("Use current web search");
    expect(prompt).toContain("Provenance Appendix");
  });

  it("recognizes supported modes", () => {
    expect(isMode("discover")).toBe(true);
    expect(isMode("unknown")).toBe(false);
  });
});
