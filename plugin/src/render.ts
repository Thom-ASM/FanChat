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
  const menu = document.createElement("div");
  menu.id = "fan-chat-stream-menu";

  getRoot().appendChild(menu);
}

function loadStream() {}
