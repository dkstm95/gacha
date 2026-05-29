import { describe, expect, it } from "vitest";
import { renderArgs, selectPlatform } from "../src/platforms/runner.js";
import type { InvestiqConfig } from "../src/types.js";

describe("platform runner", () => {
  it("renders prompt placeholders in args", () => {
    expect(renderArgs(["-p", "{{prompt}}"], "hello")).toEqual(["-p", "hello"]);
  });

  it("falls back to manual when no platform is enabled", () => {
    const config: InvestiqConfig = {
      version: 1,
      defaultPlatform: "auto",
      platformPriority: ["missing", "manual"],
      requireFreshData: true,
      allowTradeExecution: false,
      platforms: {
        missing: {
          label: "Missing",
          command: "definitely-not-installed-investiq-test",
          args: ["{{prompt}}"],
          promptMode: "argument",
          subscription: "",
          enabled: true
        },
        manual: {
          label: "Manual",
          command: "",
          args: [],
          promptMode: "print",
          subscription: "manual",
          enabled: true
        }
      }
    };

    expect(selectPlatform(config, "auto")).toBe("manual");
  });
});
