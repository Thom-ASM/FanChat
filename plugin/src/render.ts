import { renderState } from "./types.ts";
import { createMaterialIcon } from "./utils.ts";

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
        tearDown();
        popup();
        break;

      case "select-stream":
        tearDown();
        loadStream();
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

  const text = document.createElement("p");
  text.innerText = "Chat";

  button.append(text);

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

  menu.append(navBar("open-menu"), streams);

  const streamNodes = streamData.events.map((liveStream) => {
    return streamContainer(liveStream);
  });

  streams.append(...streamNodes);

  getRoot().appendChild(menu);
}

function loadStream() {
  const frame = document.createElement("iframe");

  frame.referrerPolicy = "origin";
  frame.src =
    "https://www.youtube.com/live_chat" +
    "?v=Dx5qFachd3A" +
    "&embed_domain=play.nugs.net" +
    "&dark_theme=1";

  frame.style.position = "fixed";
  frame.style.top = "20px";
  frame.style.right = "20px";
  frame.style.width = "400px";
  frame.style.height = "700px";
  frame.style.zIndex = "2147483647";

  getRoot().append(navBar(), frame);
}

function livePill() {
  const pill = document.createElement("span");

  pill.id = "fan-chat-live-pill";

  pill.innerText = "LIVE";

  return pill;
}

function streamContainer(liveStream: stream) {
  const streamContainer = document.createElement("div");

  streamContainer.id = "fan-chat-stream-container";

  const streamContent = document.createElement("div");
  const streamTitle = document.createElement("h2");
  streamTitle.innerText = liveStream.artist;

  const pill = livePill();

  streamContent.append(streamTitle, pill);

  const loadChat = document.createElement("button");
  loadChat.innerText = "Load Chat";

  loadChat.addEventListener("click", () => {
    loadChat.dispatchEvent(
      new CustomEvent("fan-chat-action", {
        bubbles: true,
        detail: {
          action: "select-stream",
        },
      }),
    );
  });

  streamContainer.append(streamContent, loadChat);

  return streamContainer;
}

function navBar(state: "open-menu" | "close-menu" | "select-stream") {
  const nav = document.createElement("div");
  nav.id = "fan-chat-nav-bar";

  let nodes = [];

  const backButton = document.createElement("button");
  backButton.id = "fan-chat-nav-bar-back-button";
  backButton.innerText = state === "open-menu" ? "Close" : "Back";

  backButton.addEventListener("click", () => {
    backButton.dispatchEvent(
      new CustomEvent("fan-chat-action", {
        bubbles: true,
        detail: {
          action: "close-menu",
        },
      }),
    );
  });

  nodes.push(backButton);

  if (state === "open-menu") {
    const title = document.createElement("h1");
    title.id = "fan-Chat-Title";
    title.innerText = "Fan Chat";

    nodes.push(title);
  }

  nav.append(...nodes);

  return nav;
}
