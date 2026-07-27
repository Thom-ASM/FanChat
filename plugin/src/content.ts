import { stream } from "./types.ts";
import { render } from "./render.ts";
import "./styles.css";

function fetchStreams(): Promise<Array<stream>> {
  return [
    { artisName: "Billy Strings", startTime: "1785181699", ytChatURL: "" },
  ];
}

function init() {
  const streams = fetchStreams();

  render();
}

init();
