// Renders each content/blogs/*/diagrams/*.mmd to an .svg beside it, for the
// diagram shortcode to inline. Mermaid sizes every box from its label's
// measured width, so Inter (the site's text font) is loaded before anything
// renders; otherwise labels are measured in a fallback font and clip once the
// page's Inter applies.
//
//   node tools/diagrams/render.mjs [file.mmd ...]
import { readFile, writeFile, glob } from "node:fs/promises";
import { createRequire } from "node:module";
import path from "node:path";
import puppeteer from "puppeteer";

const root = path.resolve(import.meta.dirname, "../..");
const require = createRequire(import.meta.url);
const config = JSON.parse(await readFile(path.join(import.meta.dirname, "config.json"), "utf8"));
const font = (await readFile(path.join(root, "assets/fonts/inter-latin-wght-normal.woff2"))).toString("base64");

let files = process.argv.slice(2);
if (!files.length) files = await Array.fromAsync(glob("content/blogs/*/diagrams/*.mmd", { cwd: root }));

const browser = await puppeteer.launch();
const page = await browser.newPage();
await page.setContent(`<!doctype html><style>
  @font-face { font-family: "Inter var"; src: url(data:font/woff2;base64,${font}) format("woff2"); font-weight: 100 900; }
  body { line-height: 1.5; }
</style><body></body>`);
await page.addScriptTag({ path: require.resolve("mermaid/dist/mermaid.min.js") });
await page.evaluate(async (config) => {
  await document.fonts.load('15px "Inter var"');
  mermaid.initialize({ ...config, startOnLoad: false });
}, config);

for (const f of files.sort()) {
  const src = path.resolve(root, f);
  const id = "mmd-" + path.basename(src, ".mmd");
  const code = await readFile(src, "utf8");
  const svg = await page.evaluate(async (id, code) => (await mermaid.render(id, code)).svg, id, code);
  await writeFile(src.replace(/\.mmd$/, ".svg"), svg + "\n");
  console.log("rendered", path.relative(root, src.replace(/\.mmd$/, ".svg")));
}
await browser.close();
