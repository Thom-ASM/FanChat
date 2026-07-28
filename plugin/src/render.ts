import { renderState } from "./types.ts";

export function render() {
  const root = document.createElement("div");
  root.id = "fan-chat-root";

  root.addEventListener("fan-chat-action", (event) => {
    const customEvent = event as CustomEvent<{
      action: "open-menu" | "close-menu" | "select-stream";
      streamId?: string;
    }>;

    switch (customEvent.detail.action) {
      case "open-menu":
        tearDown();
        menu();
        break;
      case "close-menu":
        // Update state, then close the panel.
        break;

      case "select-stream":
        // Use customEvent.detail.streamId.
        break;
    }
  });

  // create inital event
  //

  document.body.append(root);
  // inital state is popup
  popup();
}

function tearDown() {
  const root = getRoot();
  root.removeChild(root.children[0]);
}

function getRoot() {
  return document.getElementById("fan-chat-root");
}

function popup() {
  const button = document.createElement("div");
  button.id = "fan-chat-popup";

  button.addEventListener("click", () => {
    button.dispatchEvent(
      new CustomEvent("fan-chat-action", {
        bubbles: true,
        detail: {
          action: "open-menu",
        },
      }),
    );
  });

  getRoot().appendChild(button);
}

function menu() {
  const streamData = JSON.parse(
    localStorage.getItem("fan-chat-stream-manifest"),
  );
  const menu = document.createElement("div");
  menu.id = "fan-chat-stream-menu";

  const streams = document.createElement("div");
  streams.id = "fan-chat-streams-list";

  menu.appendChild(streams);

  const streamNodes = streamData.events.map((liveStream) => {
    const streamNode = document.createElement("button");

    streamNode.id = `fan-chat-stream-${liveStream.artist}`;

    streamNode.innerText = liveStream.artist;
    return streamNode;
  });

  streams.append(...streamNodes);

  getRoot().appendChild(menu);
}

function loadStream() {
  const frame = document.createElement("iframe");

  frame.referrerPolicy = "origin";
  frame.src =
    "https://www.youtube.com/live_chat" +
    "?v=rJt1bdqxSn0" +
    "&embed_domain=play.nugs.net";

  frame.style.position = "fixed";
  frame.style.top = "20px";
  frame.style.right = "20px";
  frame.style.width = "400px";
  frame.style.height = "700px";
  frame.style.zIndex = "2147483647";

  document.body.appendChild(frame);
}
