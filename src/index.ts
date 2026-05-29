import { runCli } from "./cli.js";

runCli(process.argv.slice(2)).then((code) => {
  process.exit(code);
}).catch((error: unknown) => {
  console.error(error instanceof Error ? error.message : String(error));
  process.exit(1);
});
