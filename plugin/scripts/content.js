function loadChat() {


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

loadChat()
