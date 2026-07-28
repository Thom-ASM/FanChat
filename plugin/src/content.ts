import { stream } from "./types.ts";
import { render } from "./render.ts";
import "./styles.css";

async function fetchStreams(): Promise<Array<stream>> {
  return fetch("https://d1rwrc4jiryi4e.cloudfront.net/manifest.json", {
    //todo: remove
    cache: "no-store",
  });
}

async function init() {
  const streams = await fetchStreams();
  const json = await streams.json();
  console.log("json", json);

  localStorage.setItem("fan-chat-stream-manifest", JSON.stringify(json));

  render();
}

init();
