(function() {
	//#region src/render.ts
	function render() {
		const root = document.createElement("div");
		root.id = "fan-chat-root";
		root.addEventListener("fan-chat-action", (event) => {
			switch (event.detail.action) {
				case "open-menu":
					tearDown();
					menu();
					break;
				case "close-menu": break;
				case "select-stream": break;
			}
		});
		document.body.append(root);
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
			button.dispatchEvent(new CustomEvent("fan-chat-action", {
				bubbles: true,
				detail: { action: "open-menu" }
			}));
		});
		getRoot().appendChild(button);
	}
	function menu() {
		const menu = document.createElement("div");
		menu.id = "fan-chat-stream-menu";
		getRoot().appendChild(menu);
	}
	//#endregion
	//#region src/content.ts
	function init() {
		render();
	}
	init();
	//#endregion
})();

//# sourceMappingURL=content.js.map